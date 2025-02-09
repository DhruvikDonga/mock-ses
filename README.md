# Mock SES

## Mock Send Email API

Mock API takes the [request structure](https://docs.aws.amazon.com/ses/latest/APIReference-V2/API_SendEmail.html#API_SendEmail_RequestSyntax) as of Send Email V2 docs and it saves the random status of that email as of what [SES saves](https://docs.aws.amazon.com/ses/latest/dg/monitor-sending-activity.html#:~:text=event%20types%20to-,monitor%20in%20SES%3A,-Send%20%E2%80%93%20The%20send).

Note:- So in request structure the mock API only sends simple message in content there are 3 options though so follow this doc if you want to create [content simple](https://docs.aws.amazon.com/ses/latest/APIReference-V2/API_EmailContent.html)

Mock API takes request headers:-
Check [AWS signature v4](https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_sigv.html) 
[Sample implementation](https://gist.github.com/anandkunal/b67eb94454b77cfc2b50026989586cc0#file-aws_sigv4_ses-go-L101)
- Authorization (simply checks prefix :- "AWS4-HMAC-SHA256") 
- X-Amz-Date

### Mock SES Email Sending Rules
AWS SES Rules hardcode implementation in mock API [Reference](https://docs.aws.amazon.com/ses/latest/dg/manage-sending-quotas-request-increase.html)

Email Warming Up

- New senders have lower limits (e.g., max 200 emails/day initially).
  - Limits increase gradually based on reputation.
  - Error: ThrottlingException (400) – "Email sending limit exceeded for new senders."

- Restrict emails per second (EPS) to prevent spam.
  - Example: Allow 5 emails per second, delay excess requests.
  - Error: LimitExceededException (400) – "Daily email sending limit reached."

- Simulate bounced emails (5% chance).
  - Track complaint rate (must stay below 0.1%).
  - If bounce/complaints exceed limits, temporarily block the sender.
  - Error: AccountSuspendedException (403) – "Sender temporarily blocked due to high bounce rate."

- Maintain a list of suppressed emails.
  - Block sending to emails that have bounced multiple times.
  - Domain & Email Verification
  - Error: EmailAddressSuppressed (400) – "This email address is on the suppression list."

- Reject emails from unverified domains.
  - Ensure from_email exists in a verified sender list.
  - Request Throttling & Temporary Failures
  - Error: NotAuthorized (401) – "Sender email not verified."

- Assign statuses like Sent, Bounced, Delivered, Rejected, Complaint, etc., randomly for realistic testing.

What can be added :-
- The bounce algorithm should be improved cause once percentage issue comes up that sender is blocklisted need to add a feature to unblock it

## API Endpoints

1. **Ping** :- To check health of app
   ```bash 
    curl  -X GET \
    'http://localhost:8080/ping' \
    --header 'Accept: */*' \
    --header 'User-Agent: Thunder Client (https://www.thunderclient.com)'
   ```

   **Response**
   ```json
    {
      "API service": "UP",
      "DB service": "UP",
      "Timestamp": "2025-02-09T11:37:14.921742278Z"
    }  
   ```

2. **Send Mock Email** :- Mock test API
   1. Check if email is verified or add it
   2. Check if email has proper email limits or add it

    ```bash
    curl  -X POST \
    'http://localhost:8080/v2/email/outbound-emails' \
    --header 'Accept: */*' \
    --header 'User-Agent: Thunder Client (https://www.thunderclient.com)' \
    --header 'Authorization: AWS4-HMAC-SHA256 Credential=mock-credential' \
    --header 'X-Amz-Date: 20250208T120000Z' \
    --header 'Content-Type: application/json' \
    --data-raw '{
    "FromEmailAddress": "dongadhruvik@gmail.com",
    "Destination": {
      "ToAddresses": ["recipient@example.com"],
      "CcAddresses": ["cc@example.com"],
      "BccAddresses": ["bcc@example.com"]
    },
    "Content": {
      "Simple": {
        "Subject": {
          "Data": "Test Email Subject",
          "Charset": "UTF-8"
        },
        "Body": {
          "Text": {
            "Data": "This is a test email in plain text.",
            "Charset": "UTF-8"
          },
          "Html": {
            "Data": "<h1>This is a test email in HTML</h1>",
            "Charset": "UTF-8"
          }
        }
      }
    }
    }
    '
    ```

    **Response**
    ```json
    {
      "MessageId": "74319a97-2cd6-4183-bffa-dabb4b0ab204"
    }
    ```

3. **Set Email Limit** :- To set the email limit of sender
    ```bash
    curl  -X POST \
    'http://localhost:8080/email-limits' \
    --header 'Accept: */*' \
    --header 'User-Agent: Thunder Client (https://www.thunderclient.com)' \
    --header 'Content-Type: application/json' \
    --data-raw '{
      "sender_email": "dongadhruvik@gmail.com",
      "daily_quota": 2
      }
    '
    ```
  
    **Response**
    ```json
    {
      "message": "Email limits updated successfully"
    }
    ```
4. **Add/Delete Email to Suppressed List** :- Recipient emails to suppressed list
     change POST to DELETE to remove email 
    ```bash
    curl  -X POST \
    'http://localhost:8080/suppression-list' \
    --header 'Accept: */*' \
    --header 'User-Agent: Thunder Client (https://www.thunderclient.com)' \
    --header 'Content-Type: application/json' \
    --data-raw '{
      "email": "recipient@example.com"
    }
    '
    ```

6. **Add/Delete Verified Email Sender List**
    change POST to DELETE to remove email 
   ```bash
    curl  -X POST \
    'http://localhost:8080/verified-senders' \
    --header 'Accept: */*' \
    --header 'User-Agent: Thunder Client (https://www.thunderclient.com)' \
    --header 'Content-Type: application/json' \
    --data-raw '{
      "email": "dongadhruvik@gmail.com"
    }
    '
    ```
7. Message ID stats

   ```bash
   curl  -X GET \
    'http://localhost:8080/email/stats/f5220e28-b0e8-456a-b96b-9779352bded4' \
    --header 'Accept: */*' \
    --header 'User-Agent: Thunder Client (https://www.thunderclient.com)'
   ``` 

   Response
   ```json
   {
    "logs": {
      "Logs": [
      {
        "recipient_email": "recipient@example.com",
        "sender_email": "dongadhruvik@gmail.com",
        "status": "Click",
        "response": "Recipient clicked on a link in the email.",
        "created_at": "2025-02-09T14:28:55.614179Z"
      },
      {
        "recipient_email": "cc@example.com",
        "sender_email": "dongadhruvik@gmail.com",
        "status": "Subscription",
        "response": "Recipient has unsubscribed.",
        "created_at": "2025-02-09T14:28:55.614179Z"
      },
      {
        "recipient_email": "bcc@example.com",
        "sender_email": "dongadhruvik@gmail.com",
        "status": "Complaint",
        "response": "Recipient marked email as spam.",
        "created_at": "2025-02-09T14:28:55.614179Z"
      }
    ]
    },
    "message_id": "f5220e28-b0e8-456a-b96b-9779352bded4"
    }
   ```
  

## Tasks
- [x] Project Setup
  - [x] Go App setup
  - [x] Database (Postgres Docker container)
- [x] Database Structure
  - [x] emails table
  - [x] email recipients table
  - [x] email delivery logs table
  - [x] email statistics table
- [x] [Mock Send Email v2](https://docs.aws.amazon.com/ses/latest/APIReference-V2/API_SendEmail.html) implementation
  - [x] Verified Email API 
  - [x] Suppressed Email API
  - [x] Mock API email 
  - [x] Set Email Limits
- [x] Email Delivery Stats API

## Useful Commands
- ```docker compose build app```
- ```docker compose build```
- ```docker compose up```
- ```docker-compose exec db psql -U postgres -d mock_ses```