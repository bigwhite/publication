package context

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Skill struct {
	Name        string
	Description string
	Body        string
}

type SkillLoader struct {
	workDir string
}

func NewSkillLoader(workDir string) *SkillLoader {
	return &SkillLoader{workDir: workDir}
}

func (s *SkillLoader) LoadAll() string {
	skillBaseDir := filepath.Join(s.workDir, ".claw", "skills")

	if _, err := os.Stat(skillBaseDir); os.IsNotExist(err) {
		return ""
	}

	var skillsBuilder strings.Builder
	skillsBuilder.WriteString("\n### 可用专业技能 (Agent Skills)\n")
	skillsBuilder.WriteString("以下是你拥有的标准化外挂技能，请在符合 description 描述的场景下严格遵循其正文指令：\n\n")

	err := filepath.WalkDir(skillBaseDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && d.Name() == "SKILL.md" {
			content, err := os.ReadFile(path)
			if err == nil {
				skill := parseSkillMD(string(content))

				skillsBuilder.WriteString(fmt.Sprintf("#### 技能名称: %s\n", skill.Name))
				skillsBuilder.WriteString(fmt.Sprintf("**触发条件**: %s\n\n", skill.Description))
				skillsBuilder.WriteString("**执行指南**:\n")
				skillsBuilder.WriteString(skill.Body)
				skillsBuilder.WriteString("\n\n---\n")
			}
		}
		return nil
	})

	if err != nil || skillsBuilder.Len() < 50 {
		return ""
	}

	return skillsBuilder.String()
}

func parseSkillMD(content string) Skill {
	skill := Skill{
		Name:        "Unknown Skill",
		Description: "No description provided.",
		Body:        content,
	}

	if strings.HasPrefix(content, "---\n") || strings.HasPrefix(content, "---\r\n") {
		parts := strings.SplitN(content, "---", 3)
		if len(parts) == 3 {
			frontmatter := parts[1]
			skill.Body = strings.TrimSpace(parts[2])

			lines := strings.Split(frontmatter, "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "name:") {
					skill.Name = strings.TrimSpace(strings.TrimPrefix(line, "name:"))
				} else if strings.HasPrefix(line, "description:") {
					skill.Description = strings.TrimSpace(strings.TrimPrefix(line, "description:"))
				}
			}
		}
	}

	return skill
}
