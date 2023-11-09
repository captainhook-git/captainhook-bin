package types

type Commit struct {
	Hash    string
	Names   string
	Date    string
	Author  string
	Subject string
	Body    string
}

func CreateCommitFromXML(xml *XmlCommit) *Commit {
	c := Commit{
		Hash:    xml.Hash,
		Names:   xml.Names,
		Date:    xml.Date,
		Author:  xml.Author,
		Subject: xml.Subject,
		Body:    xml.Body,
	}
	return &c
}
