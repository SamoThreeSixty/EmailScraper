# Email Scraper

The task of this is to make a flexable microservice that can connect via IMAP to an email service, and periodically scrape for new emails.
The contents will be downloaded to sql, and the attachments will be saved to local files if they are not a security threat.

This service can then be used within other applications to provide an interface where emails can be displayed.