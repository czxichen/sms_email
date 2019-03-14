package sns

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

// SendCode 发送验证码
func SendCode(cfg *SendCodeConfig) (uint, error) {
	send, exist := SenderGet(cfg.Provider, cfg.Type)
	if !exist {
		return 0, fmt.Errorf("unsupport %s%s provider", cfg.Provider, cfg.Type)
	}

	db := BeginDB()
	defer db.DBRollback()

	var codeInfo = &CodeInfo{Application: cfg.Application, Tag: cfg.Tag, Contact: cfg.Contact}
	gdb := db.Where("try_num < max_try").Order("id DESC").Limit(1).Find(codeInfo, codeInfo)
	if gdb.Error != nil && !gdb.RecordNotFound() {
		return 0, gdb.Error
	}

	var now = time.Now()
	if gdb.RecordNotFound() || now.After(codeInfo.EnableSame) {
		codeInfo.CreatedAt = now
		codeInfo.UpdatedAt = now
		codeInfo.Expiry = now.Add(time.Duration(cfg.Expired) * time.Second)
		if cfg.EnableSame == 0 {
			cfg.EnableSame = cfg.Expired
		}
		codeInfo.EnableSame = now.Add(time.Duration(cfg.EnableSame) * time.Second)
		codeInfo.ID = 0
		codeInfo.Code = randCode(6)
		codeInfo.TryNum = 0
		codeInfo.MaxTry = cfg.MaxTry

		if err := db.Create(codeInfo).Error; err != nil {
			return 0, err
		}
	}

	db.DBCommit()
	if cfg.Data == nil {
		cfg.Data = make(map[string]interface{})
	}
	cfg.Data["Code"] = codeInfo.Code
	extend, content, err := TemplateToMsg(cfg)
	if err != nil {
		return 0, err
	}
	err = send.Send(cfg.Contact, extend, content)
	return codeInfo.ID, err
}

// CheckCode 校验验证码
func CheckCode(cfg *CheckCodeConfig) (bool, error) {
	db := BeginDB()
	defer db.DBRollback()

	var codeInfo CodeInfo
	if gdb := db.Where("id = ? and (expiry >= ? and try_num <= max_try)", cfg.ID, time.Now()).Find(&codeInfo); gdb.Error != nil {
		if gdb.RecordNotFound() {
			return false, nil
		}
		return false, gdb.Error
	}

	codeValid := codeInfo.Code == cfg.Code
	if codeValid {
		// 删除失效的验证码记录
		db.Where("expiry < ? OR try_num > max_try OR id = ?", time.Now(), cfg.ID).Delete(&CodeInfo{})
	} else {
		// 更新尝试次数
		if err := db.Model(&codeInfo).Update("try_num", gorm.Expr("try_num + 1")).Error; err != nil {
			return false, err
		}
	}
	db.DBCommit()
	return codeValid, nil
}
