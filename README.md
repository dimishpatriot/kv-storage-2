# kv-storage-2

learning project key-value storage service with local file storage & postgres storage

## config
- rename `.env-backbone` file to `.env` and fill in the required data. 
  
  **! do not add env-file to the repository !**
- configure and run the Postgres server if necessary

## run
`go run . -s=<type-of-storage>`, where type is:
- `local` - local file storage
- `postgres` - postgres storage