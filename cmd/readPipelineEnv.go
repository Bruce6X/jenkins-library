package cmd

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/SAP/jenkins-library/pkg/config"
	"github.com/SAP/jenkins-library/pkg/log"
	"github.com/SAP/jenkins-library/pkg/orchestrator"
	"github.com/SAP/jenkins-library/pkg/piperenv"
	"github.com/spf13/cobra"
	"io"
	"os"
	"path"
)

type readPipelineEnvOptions struct {
	Secret string `json:"secret,omitempty"`
}

// ReadPipelineEnv reads the commonPipelineEnvironment from disk and outputs it as JSON
func ReadPipelineEnv() *cobra.Command {
	const STEP_NAME = "readPipelineEnv"
	var stepConfig readPipelineEnvOptions
	metadata := readPipelineEnvMetadata()

	return &cobra.Command{
		Use:   "readPipelineEnv",
		Short: "Reads the commonPipelineEnvironment from disk and outputs it as JSON",
		PreRun: func(cmd *cobra.Command, args []string) {
			path, _ := os.Getwd()
			fatalHook := &log.FatalHook{CorrelationID: GeneralConfig.CorrelationID, Path: path}
			log.RegisterHook(fatalHook)

			err := PrepareConfig(cmd, &metadata, STEP_NAME, &stepConfig, config.OpenPiperFile)
			if err != nil {
				log.SetErrorCategory(log.ErrorConfiguration)
				return
			}
			log.RegisterSecret(stepConfig.Secret)
		},

		Run: func(cmd *cobra.Command, args []string) {
			err := runReadPipelineEnv(&stepConfig)
			if err != nil {
				log.Entry().Fatalf("error when writing reading Pipeline environment: %v", err)
			}
		},
	}
}

func runReadPipelineEnv(config *readPipelineEnvOptions) error {
	cpe := piperenv.CPEMap{}

	err := cpe.LoadFromDisk(path.Join(GeneralConfig.EnvRootPath, "commonPipelineEnvironment"))
	if err != nil {
		return err
	}

	// try to encrypt
	if config.Secret != "" && orchestrator.DetectOrchestrator() != orchestrator.Jenkins {
		log.Entry().Debug("found PIPER_pipelineEnv_SECRET, trying to encrypt CPE")
		jsonBytes, _ := json.Marshal(cpe)
		encrypted, err := encrypt([]byte(config.Secret), jsonBytes)
		if err != nil {
			log.Entry().Fatal(err)
		}

		os.Stdout.Write(encrypted)
		return nil
	}

	// fallback
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "\t")
	if err := encoder.Encode(cpe); err != nil {
		return err
	}

	return nil
}

func encrypt(secret, inBytes []byte) ([]byte, error) {
	// use SHA256 as key
	key := sha256.Sum256(secret)
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, fmt.Errorf("failed to create new cipher: %v", err)
	}

	// Make the cipher text a byte array of size BlockSize + the length of the message
	cipherText := make([]byte, aes.BlockSize+len(inBytes))

	// iv is the ciphertext up to the blocksize (16)
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, fmt.Errorf("failed to init iv: %v", err)
	}

	// Encrypt the data:
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], inBytes)

	// Return string encoded in base64
	return []byte(base64.StdEncoding.EncodeToString(cipherText)), err
}

// retrieve step metadata
func readPipelineEnvMetadata() config.StepData {
	var theMetaData = config.StepData{
		Metadata: config.StepMetadata{
			Name: "readPipelineEnvMetadata",
		},
		Spec: config.StepSpec{
			Inputs: config.StepInputs{
				Parameters: []config.StepParameters{
					{
						Name: "secret",
						ResourceRef: []config.ResourceReference{
							{
								Name: "cpeSecret",
								Type: "vaultSecret",
							},
						},
						Type:      "string",
						Mandatory: false,
						Default:   os.Getenv("PIPER_pipelineEnv_SECRET"),
					},
				},
			},
		},
	}
	return theMetaData
}
