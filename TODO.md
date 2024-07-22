- onboarding providers: AWS, GCP, Azure
- viewing resources: AWS, GCP, Azure
- reset password
- logout
- account details
- modify primary admin for organization
- resend email verification on the verification screen
- create another admin
- modify existing admins, including role
- enforce password complexity
- non-verified admin flows: login failure prompts for email verification resend AND signup for existing org with existing primary admin email prompts for email verification resend


```
insert into providers(external_identifier, display_name, provider_name, metadata, organization_id) values('ACCOUNT_ID', 'My AWS acct', 'AWS', '{"awsMetadata": {"role_arn": "arn:aws:iam::ACCOUNT_ID:role/ReadOnlyRole"}}', '51207fbd-87dd-48bb-b9b8-904832ead230');
```