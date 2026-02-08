# Ansible Deployment Guide

This guide explains how to deploy the Simple Go API application to an Ubuntu server using Ansible.

## Prerequisites

1.  **Ansible Installed:** You need Ansible installed on your local machine.
    *   **Mac:** `brew install ansible`
    *   **Linux:** `sudo apt install ansible`
2.  **Server Access:** You must have `root` SSH access to your Ubuntu server.

## 1. Configure Inventory

Open `ansible/inventory.ini` and update the connection details:

```ini
[webservers]
# Replace x.x.x.x with your server IP
production_server ansible_host=x.x.x.x ansible_user=root ansible_ssh_private_key_file=~/.ssh/id_rsa
```

## 2. Deploy to Server

Run the playbook from the project root directory. You must provide the sensitive configuration values using the `--extra-vars` flag.

### Command Format

```bash
ansible-playbook -i ansible/inventory.ini ansible/playbook.yml \
  --extra-vars "docker_username=machacek.j" \
  --extra-vars "db_root_password=YOUR_SECURE_ROOT_PASSWORD" \
  --extra-vars "db_user=api_user" \
  --extra-vars "db_password=YOUR_SECURE_DB_PASSWORD" \
  --extra-vars "jwt_secret=YOUR_LONG_RANDOM_SECRET_STRING" \
  --extra-vars "app_url=http://YOUR_SERVER_IP"
```

### What This Does
1.  **System Setup:** Installs Docker, Docker Compose, and necessary system tools.
2.  **Security:** Enables UFW firewall and allows ports 22 (SSH), 8080 (API), and 3000 (Frontend).
3.  **Deployment:**
    *   Creates `/opt/simple-go-api`.
    *   Generates `docker-compose.yml` and `.env` with your secrets.
    *   Pulls the latest images from Docker Hub.
    *   Starts the application.

## Troubleshooting

*   **SSH Error:** Ensure your SSH key path in `inventory.ini` is correct and has the right permissions (`chmod 600 key.pem`).
*   **Frontend Connection:** If the frontend cannot reach the backend, ensure `app_url` matches exactly what you type in the browser (e.g., `http://1.2.3.4`).
