# 前端构建层，构建前端代码，生成静态文件
FROM node:10 AS build-ui

# 拷贝项目代码
ADD . /root/ui
# 设置工作目录
WORKDIR /root/ui

# 安装依赖
RUN npm install
# 构建，生成静态文件
RUN npm run build


# 运行层
FROM nginx AS run
# 拷贝上一层构建后的静态文件目录到当前层
COPY --from=build-ui /root/ui/build /usr/share/nginx/html
EXPOSE 80

# 容器入口，启动容器时运行该命令，且不会被 docker run 提供的命令覆盖
ENTRYPOINT ["nginx", "-g", "daemon off;"]
