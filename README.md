# github-backup

Backup your GitHub repositories (including issues and comments).

The backup will include a copy of the git repository itself (a bare, mirror
clone), the wiki repository if there is one, all of the repository issues and
comments (including the important metadata for each issue and comment and any
attached files), and all of the pull requests (including TBD). All information
will be written to disk for easy follow-up with e.g., `tar` or similar and
copy to the backup medium.
