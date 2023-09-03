#### How to generate a modified jenkins-library
1. Modify golang code in ./cmd
2. Change `SAP/jenkins-library` to ${GITHUB_REPOSITORY} in https://github.com/Bruce6X/jenkins-library/blob/master/.github/workflows/upload-go-master.yml
3. Change `SAP/jenkins-library` to ${GITHUB_REPOSITORY} in https://github.com/Bruce6X/jenkins-library/blob/master/.github/workflows/release-go.yml
4. To be continued: may refer to https://github.com/HoffmannThomas/jenkins-library/commits/v1.71.4
