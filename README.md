# Mock SES

## Mock Send Email api

Mock API takes the [request structure](https://docs.aws.amazon.com/ses/latest/APIReference-V2/API_SendEmail.html#API_SendEmail_RequestSyntax) as of Send Email V2 docs and it saves the random status of that email as of what [SES saves](https://docs.aws.amazon.com/ses/latest/dg/monitor-sending-activity.html#:~:text=event%20types%20to-,monitor%20in%20SES%3A,-Send%20%E2%80%93%20The%20send) .

Note:- So in request structrue the mock api only sends simple message in content there are 3 options though so follow this doc if you want to create [content simple](https://docs.aws.amazon.com/ses/latest/APIReference-V2/API_EmailContent.html)

Mock API takes request headers :-
Check [AWS signatuer v4](https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_sigv.html) 
[Sample implementation](https://gist.github.com/anandkunal/b67eb94454b77cfc2b50026989586cc0#file-aws_sigv4_ses-go-L101)
- Authorization (simply checks prefix :- "AWS4-HMAC-SHA256") 
- X-Amz-Date

### Mock SES email sending rules
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

## Tasks
- [x] Project Setup
  - [x] Go App setup
  - [x] Database (Postgres Docker container)
- [x] Database Structure
  - [x] emails table
  - [x] email receipents table
  - [x] email delivery logs table
  - [x] email statistics table
- [] [Mock Send Email v2](https://docs.aws.amazon.com/ses/latest/APIReference-V2/API_SendEmail.html) implementation