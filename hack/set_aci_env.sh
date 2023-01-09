#!/usr/bin/env bash
# 功能：
# 该脚本用于初始化 aci Job 的构建环境，需要和 antcode 配合使用

# 根据 秘钥(kusion_aci_ssh_key) 生成 ssh 私钥和公钥, 并设置权限
# NOTE: 变量 kusion_aci_ssh_key 需要在仓库中配置
# ssh-key 私钥对应的环境变量名和每个仓库配置的公钥key名字保持一致
# ssh-key 对应的环境变量不包含 BEGIN RSA PRIVATE KEY 和 "END RSA PRIVATE KEY" 部分
# https://aci.antgroup.com/#/project/

echo "-----BEGIN RSA PRIVATE KEY-----" > ~/.ssh/id_rsa
echo  ${ACI_VAR_kusion_aci_ssh_key} |tr -s '[:space:]' '\n' >> ~/.ssh/id_rsa
echo "-----END RSA PRIVATE KEY-----" >> ~/.ssh/id_rsa

chmod 600 ~/.ssh/id_rsa
chmod 700 ~/.ssh

ssh-keygen -y -f ~/.ssh/id_rsa > ~/.ssh/id_rsa.pub
chmod 600 ~/.ssh/id_rsa.pub

echo "StrictHostKeyChecking no" >> ~/.ssh/config
chmod 600 ~/.ssh/config

# 添加可信的网址的公钥

ssh-keyscan -t rsa github.com >> ~/.ssh/known_hosts
ssh-keyscan -t rsa code.alipay.com >> ~/.ssh/known_hosts
ssh-keyscan -t rsa gitlab.alipay-inc.com >> ~/.ssh/known_hosts

# Go: 配置私有仓库访问方式
git config --global url."git@code.alipay.com:".insteadOf "https://code.alipay.com"
git config --global url."git@gitlab.alipay-inc.com:".insteadOf "http://gitlab.alipay-inc.com"
git config --global url."git@gitlab.alibaba-inc.com:".insteadOf "http://gitlab.alibaba-inc.com"
git config --global url."ssh://git@code.alipay.com".insteadOf "https://code.alipay.com"
git config --global url."ssh://git@gitlab.alipay-inc.com".insteadOf "https://gitlab.alipay-inc.com"
git config --global url."ssh://git@gitlab.alibaba-inc.com".insteadOf "https://gitlab.alibaba-inc.com"

# Go 关闭私有仓库的 sum 验证策略, 并配置代理
go env -w GOPRIVATE="*.alipay-inc.com,*.alibaba-inc.com,*.alipay.com"
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GONOSUMDB=*

# Go 环境就绪, 进入工作目录进行单元测试
