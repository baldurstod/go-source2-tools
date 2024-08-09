package particles

import "github.com/baldurstod/go-source2-tools/kv3"

type Driver struct {
	ControlPoint   int
	AttachType     string
	AttachmentName string
	EntityName     string
}

func (driver *Driver) initFromDatas(datas *kv3.Kv3Element) error {
	var ok bool
	driver.ControlPoint, _ = datas.GetIntAttribute("m_iControlPoint")
	if driver.AttachType, ok = datas.GetStringAttribute("m_iAttachType"); !ok {
		driver.AttachType = "PATTACH_ABSORIGIN_FOLLOW"
	}
	driver.AttachmentName, _ = datas.GetStringAttribute("m_attachmentName")
	if driver.EntityName, ok = datas.GetStringAttribute("m_entityName"); !ok {
		driver.EntityName = "self"
	}
	return nil
}
