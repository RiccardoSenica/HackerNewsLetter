# hackernewsletter

Command to build the executable for manual upload to Lambda:
GOOS=linux CGO_ENABLED=0 go build -o main .

Required variables:

- FETCH_SIZE: how many news to fetch
- BATCH_SIZE: batch size for DB write
- REGION: AWS region
- SENDER: sender email address
- RECIPIENT: recipient email address
