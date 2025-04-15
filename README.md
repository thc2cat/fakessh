# project: fake SSH Connection Logger

This project involves creating a custom SSH server that captures and logs connection details (username, password, source IP address) to a syslog server.

This provides a simple monitoring and auditing tool for SSH access, enabling administrators to track login attempts and identify suspicious activity.

## Key Features

Accepts incoming SSH connections.
Extracts and records login credentials (username, password) and source IP.
Forwards connection information to a syslog system.
Offers flexible configuration for syslog management.

## Objectives

Develop a straightforward and effective SSH connection logging solution.

Assist administrators in monitoring and auditing access attempts.

Facilitate the detection of suspicious behavior and intrusion attempts.

## Important Notes

Logging passwords poses **significant security risks**.

Implement robust security measures to protect this sensitive data.

This project is intended for auditing and monitoring purposes.

Adhere to all applicable laws and regulations concerning personal data collection and storage.
