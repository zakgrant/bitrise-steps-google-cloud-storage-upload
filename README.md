## Bitrise GCS Upload Step

Bitrise Step for uploading build artifacts to Google Cloud Storage.

### How to use it.

#### Upload Service Account JSON File
Bitrise GCS Upload Step uses a service account json file for authentication.

Don't have a service account yet then create one -> [instructions](https://cloud.google.com/docs/authentication).

JSON Credentials file will need uploaded to Bitrise -> `Workflow -> Code Signing -> Generic File Storage`

**WARNING** Bitrise GCS Upload Step uses `GCS_SERVICE_ACCOUNT_JSON_KEY_URL` as default key. 
If you choose to use a different name, you will have to update it in step inputs.
 
#### Update `bitrise.yml`
Add following step into your `bitrise.yml`

```yaml
- git::https://github.com/zakgrant/bitrise-steps-google-cloud-storage-upload.git@1.0:
    title: Upload artefact to Google Cloud Storage
    inputs:
    - GCS_SERVICE_ACCOUNT_JSON_KEY_URL: BITRISE_GCS_SERVICE_ACCOUNT_JSON_KEY_URL
    - BUCKET_NAME: bucket-name
    - FOLDER_NAME: a-folder/another-folder
    - UPLOAD_FILE_PATH: $BITRISE_BUILD_ARTEFACT_PATH
    - UPLOADED_FILE_NAME: $GOOGLE_STORAGE_UPLOADED_FILE_NAME
```
