# Supernote Cloud Sync

An open-source Golang command-line application to sync your Supernote cloud files to a local directory for backup, version control, and offline access.

## Current Features:

- **Open source**: Contribute to the project on GitHub!
- **Login using email**: Burrently only email is supported but I can look into using mobile number
- **Get user info**: This is just an endpoint used on the website, good way to check who is logged in
- **Get equipment binding status**: The api is only available to accounts that have been connected to a device
- **Get list of files/folders**: This will be the main functionality used within the sync, also produces MD5 hash which will be helpful for delta changes. 

## Proposed Features:
- **Download files**: from the cloud.
- **Upload files**: to the cloud.
- **Two-way sync**: Download new files from the cloud and upload changes made locally.
- **Deal with file conflicts**: When a file conflict occurs, the user will be prompted to resolve the conflict.
- **Customizable**: Choose which notebooks and folders to sync, and how often.

## Disclaimer:
This project is not affiliated with Supernote. Use at your own risk.

## Motivation:
I built this tool because I wanted a simple way to backup my Supernote notes and have offline access.
Feedback and contributions welcome!

I'm still actively developing supernote-sync and would love your feedback and contributions. If you encounter any issues or have suggestions for new features, please open an issue on GitHub. Pull requests are also welcome!