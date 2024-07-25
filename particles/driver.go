package particles

import "github.com/baldurstod/go-source2-tools/kv3"

type Driver struct {
	ControlPoint   int
	AttachType     string
	AttachmentName string
	EntityName     string
}

func (driver *Driver) initFromDatas(datas *kv3.Kv3Element) error {
	driver.ControlPoint, _ = datas.GetIntAttribute("m_iControlPoint")
	driver.AttachType, _ = datas.GetStringAttribute("m_iAttachType")
	driver.AttachmentName, _ = datas.GetStringAttribute("m_attachmentName")
	driver.EntityName, _ = datas.GetStringAttribute("m_entityName")
	return nil
}
