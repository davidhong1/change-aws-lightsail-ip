# Change AWS Lightsail IP

Change aws lightsail ip when ip can't telnet in CN.

# Install

```shell
mkdir /etc/change-aws-lightsail-ip

wget https://github.com/davidhong1/change-aws-lightsail-ip/releases/download/alpha-v0.7/change-aws-lightsail-ip -O /etc/change-aws-lightsail-ip/change-aws-lightsail-ip && chmod +x /etc/change-aws-lightsail-ip/change-aws-lightsail-ip

cat > /etc/change-aws-lightsail-ip/config.yaml << EOF
aws_default_region: "ap-northeast-1" # change as you need
access_key_id: "" # change
access_secret: "" # change
default_port: 10310 # change

EOF

cat > /etc/systemd/system/change-aws-lightsail-ip.service << EOF
[Unit]
Description=Change AWS Lightsail IP
Documentation=https://github.com/davidhong1/change-aws-lightsail-ip
After=network.target nss-lookup.target

[Service]
User=nobody
NoNewPrivileges=true
ExecStart=/etc/change-aws-lightsail-ip/change-aws-lightsail-ip -config /etc/change-aws-lightsail-ip/config.yaml
Restart=on-failure
RestartPreventExitStatus=23

[Install]
WantedBy=multi-user.target

EOF

systemctl start change-aws-lightsail-ip
systemctl enable change-aws-lightsail-ip
systemctl status change-aws-lightsail-ip
```
