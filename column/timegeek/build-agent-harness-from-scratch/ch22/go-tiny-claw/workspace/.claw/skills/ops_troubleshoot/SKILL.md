---
name: ops_troubleshoot
description: Nginx 故障排查与修复标准作业程序 (SOP)。当人类报告 "服务 502"、"接口不通" 或要求排查 Nginx 错误时，必须强制加载并遵
循此技能。
---                   
                      
# Nginx 故障排查 SOP
                
你现在的角色是一线运维工程师，在排查 Nginx 故障时，请严格遵循以下排查链路：
                
1. **信息收集**：首先使用 `bash` 检查 `error.log` 的最后 50 行（例如执行：`tail -n 50 error.log`）。
2. **根因定位**：如果发现是 "upstream prematurely closed connection" 或配置文件的语法指令错误（unknown directive），请立即去检查
 `nginx.conf` 文件的具体内容。
3. **精准修复**：一旦确认配置错误，绝对不能使用 bash 的 sed 盲目替换，**必须使用 `edit_file` 工具**，提供足够上下文进行精准修正。
4. **服务重启**：修复配置后，尝试通过 `bash` 运行 `nginx -s reload` 使配置生效。系统可能会触发审批拦截，请向人类说明你重启的理由
并等待放行。

