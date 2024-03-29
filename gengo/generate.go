package gengo

import (
	"bytes"

	"text/template"
)

var msgTemplate = `
// Automatically generated from the message definition "{{ .FullName }}.msg"
package {{ .Package }}
import (
    "bytes"
{{- if .BinaryRequired }}
    "encoding/binary"
{{- end }}
{{- range .Imports }}
	"{{ . }}"
{{- end }}

    "github.com/light4d/rosgo"
)

{{- if gt (len .Constants) 0 }}
const (
{{- range .Constants }}
	{{- if eq .Type "string" }}
    {{ .GoName }} {{ .Type }} = "{{ .Value }}"
	{{- else }}
	{{ .GoName }} {{ .Type }} = {{ .Value }}
	{{- end }}
{{- end }}
)
{{- end }}


type _Msg{{ .ShortName }} struct {
    text string
    name string
    md5sum string
}

func (t *_Msg{{ .ShortName }}) Text() string {
    return t.text
}

func (t *_Msg{{ .ShortName }}) Name() string {
    return t.name
}

func (t *_Msg{{ .ShortName }}) MD5Sum() string {
    return t.md5sum
}

func (t *_Msg{{ .ShortName }}) NewMessage() rosgo.Message {
    m := new({{ .ShortName }})
{{- range .Fields }}
{{-     if .IsArray }}
{{-         if eq .ArrayLen -1 }}
	m.{{ .GoName }} = []{{ .GoType }}{}
{{-         else }}
	for i := 0; i < {{ .ArrayLen }}; i++ {
		m.{{ .GoName }}[i]  = {{ .ZeroValue }}
	}
{{-         end}}
{{-     else }}
	m.{{ .GoName }} = {{ .ZeroValue }}
{{-     end }}
{{- end }}
    return m
}

var (
    Msg{{ .ShortName }} = &_Msg{{ .ShortName }} {
        ` + "`" + `{{ .Text }}` + "`" + `,
        "{{ .FullName }}",
        "{{ .MD5Sum }}",
    }
)

type {{ .ShortName }} struct {
{{- range .Fields }}
{{-     if .IsArray }}
{{-         if eq .ArrayLen -1 }}
	{{ .GoName }} []{{ .GoType }}` + " `rosmsg:\"{{ .Name }}:{{ .Type }}[]\"`" + `
{{-         else }}
	{{ .GoName }} [{{ .ArrayLen }}]{{ .GoType }}` + " `rosmsg:\"{{ .Name }}:{{ .Type }}[{{ .ArrayLen }}]\"`" + `
{{-         end }}
{{-     else }}
	{{ .GoName }} {{ .GoType }}` + " `rosmsg:\"{{ .Name }}:{{ .Type }}\"`" + `
{{-     end }}
{{- end }}
}

func (m *{{ .ShortName }}) Type() rosgo.MessageType {
	return Msg{{ .ShortName }}
}

func (m *{{ .ShortName }}) Marshal(buf *bytes.Buffer) ([]byte,error) {
    var err error = nil
{{- range .Fields }}
{{-     if .IsArray }}
    binary.Write(buf, binary.LittleEndian, uint32(len(m.{{ .GoName }})))
    for _, e := range m.{{ .GoName }} {
{{-         if .IsBuiltin }}
{{-             if eq .Type "string" }}
        binary.Write(buf, binary.LittleEndian, uint32(len([]byte(e))))
        buf.Write([]byte(e))
{{-            else }}
{{-                if or (eq .Type "time") (eq .Type "duration") }}
        binary.Write(buf, binary.LittleEndian, e.Sec)
        binary.Write(buf, binary.LittleEndian, e.NSec)
{{-                else }}
        binary.Write(buf, binary.LittleEndian, e)
{{-                end }}
{{-             end }}
{{-         else }}
        if _,err = e.Marshal(buf); err != nil {
            return nil,err
        }
{{-         end }}
    }
{{-     else }}
{{-         if .IsBuiltin }}
{{-             if eq .Type "string" }}
    binary.Write(buf, binary.LittleEndian, uint32(len([]byte(m.{{ .GoName }}))))
    buf.Write([]byte(m.{{ .GoName }}))
{{-             else }}
{{-                 if or (eq .Type "time") (eq .Type "duration") }}
    binary.Write(buf, binary.LittleEndian, m.{{ .GoName }}.Sec)
    binary.Write(buf, binary.LittleEndian, m.{{ .GoName }}.NSec)
{{-                 else }}
    binary.Write(buf, binary.LittleEndian, m.{{ .GoName }})
{{-                 end }}
{{-             end }}
{{-         else }}
    if _,err = m.{{ .GoName }}.Marshal(buf); err != nil {
        return nil,err
    }
{{-         end }}
{{-     end }}
{{- end }}
    return buf.Bytes(),err
}


func (m *{{ .ShortName }}) Unmarshal(buf *bytes.Reader) error {
    var err error = nil
{{- range .Fields }}
{{-    if .IsArray }}
    {
        var size uint32
        if err = binary.Read(buf, binary.LittleEndian, &size); err != nil {
            return err
        }
{{-        if lt .ArrayLen 0 }}
        m.{{ .GoName }} = make([]{{ .GoType }}, int(size))
{{-        end }}
        for i := 0; i < int(size); i++ {
{{-          if .IsBuiltin }}
{{-              if eq .Type "string" }}
            {
                var size uint32
                if err = binary.Read(buf, binary.LittleEndian, &size); err != nil {
                    return err
                }
                data := make([]byte, int(size))
                if err = binary.Read(buf, binary.LittleEndian, data); err != nil {
                    return err
                }
                m.{{ .GoName }}[i] = string(data)
            }
{{-              else }}
{{- 					if or (eq .Type "time") (eq .Type "duration") }}
            {
                if err = binary.Read(buf, binary.LittleEndian, &m.{{ .GoName }}[i].Sec); err != nil {
                    return err
                }

                if err = binary.Read(buf, binary.LittleEndian, &m.{{ .GoName }}[i].NSec); err != nil {
                    return err
                }
            }
{{-                  else }}
            if err = binary.Read(buf, binary.LittleEndian, &m.{{ .GoName }}[i]); err != nil {
                return err
            }
{{-                  end }}
{{-              end }}
{{-          else }}
            if err = m.{{ .GoName }}[i].Unmarshal(buf); err != nil {
                return err
            }
{{-      	end }}
        }
    }
{{-    else }}
{{-        if .IsBuiltin }}
{{-            if eq .Type "string" }}
    {
        var size uint32
        if err = binary.Read(buf, binary.LittleEndian, &size); err != nil {
            return err
        }
        data := make([]byte, int(size))
        if err = binary.Read(buf, binary.LittleEndian, data); err != nil {
            return err
        }
        m.{{ .GoName }} = string(data)
    }
{{-            else }}
{{-            		if or (eq .Type "time") (eq .Type "duration") }}
    {
        if err = binary.Read(buf, binary.LittleEndian, &m.{{ .GoName }}.Sec); err != nil {
            return err
        }

        if err = binary.Read(buf, binary.LittleEndian, &m.{{ .GoName }}.NSec); err != nil {
            return err
        }
    }
{{-            		else }}
    if err = binary.Read(buf, binary.LittleEndian, &m.{{ .GoName }}); err != nil {
        return err
    }
{{-         			end }}
{{-            end }}
{{-        else }}
    if err = m.{{ .GoName }}.Unmarshal(buf); err != nil {
        return err
    }
{{-    	  end }}
{{-    end }}
{{- end }}
    return err
}
`

var srvTemplate = `
// Automatically generated from the message definition "{{ .FullName }}.srv"
package {{ .Package }}
import (
    "github.com/light4d/rosgo"
)

// Service type metadata
type _Srv{{ .ShortName }} struct {
    name string
    md5sum string
    text string
    reqType rosgo.MessageType
    resType rosgo.MessageType
}

func (t *_Srv{{ .ShortName }}) Name() string { return t.name }
func (t *_Srv{{ .ShortName }}) MD5Sum() string { return t.md5sum }
func (t *_Srv{{ .ShortName }}) Text() string { return t.text }
func (t *_Srv{{ .ShortName }}) RequestType() rosgo.MessageType { return t.reqType }
func (t *_Srv{{ .ShortName }}) ResponseType() rosgo.MessageType { return t.resType }
func (t *_Srv{{ .ShortName }}) NewService() rosgo.Service {
    return new({{ .ShortName }})
}

var (
    Srv{{ .ShortName }} = &_Srv{{ .ShortName }} {
        "{{ .FullName }}",
        "{{ .MD5Sum }}",
        ` + "`" + `{{ .Text }}` + "`" + `,
        Msg{{ .ShortName }}Request,
        Msg{{ .ShortName }}Response,
    }
)


type {{ .ShortName }} struct {
    Request {{ .ShortName }}Request
    Response {{ .ShortName }}Response
}

func (s *{{ .ShortName }}) ReqMessage() rosgo.Message { return &s.Request }
func (s *{{ .ShortName }}) ResMessage() rosgo.Message { return &s.Response }
`

type MsgGen struct {
	MsgSpec
	BinaryRequired bool
	Imports        []string
}

func (gen *MsgGen) analyzeImports() {
	imp_path := ""
	if len(*import_path) != 0 {
		imp_path = *import_path + "/"
	}

OUTER:
	for i, field := range gen.Fields {
		if len(field.Package) == 0 {
			gen.BinaryRequired = true
		} else if gen.Package == field.Package {
			gen.Fields[i].GoType = field.Type
			gen.Fields[i].ZeroValue = field.Type + "{}"
		} else {
			for _, imp := range gen.Imports {
				if imp == imp_path+field.Package {
					continue OUTER
				}
			}
			gen.Imports = append(gen.Imports, imp_path+field.Package)
		}

		// Binary is required to read the size of array
		if field.IsArray {
			gen.BinaryRequired = true
		}
	}
}

func GenerateMessage(context *MsgContext, spec *MsgSpec) (string, error) {
	var gen MsgGen
	gen.Fields = spec.Fields
	gen.Constants = spec.Constants
	gen.Text = spec.Text
	gen.FullName = spec.FullName
	gen.ShortName = spec.ShortName
	gen.Package = spec.Package
	gen.MD5Sum = spec.MD5Sum

	gen.analyzeImports()

	tmpl, err := template.New("msg").Parse(msgTemplate)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer

	err = tmpl.Execute(&buffer, gen)
	if err != nil {
		return "", err
	}
	return buffer.String(), err
}

func GenerateService(context *MsgContext, spec *SrvSpec) (string, string, string, error) {
	reqCode, err := GenerateMessage(context, spec.Request)
	if err != nil {
		return "", "", "", err
	}
	resCode, err := GenerateMessage(context, spec.Response)
	if err != nil {
		return "", "", "", err
	}

	tmpl, err := template.New("srv").Parse(srvTemplate)
	if err != nil {
		return "", "", "", err
	}

	var buffer bytes.Buffer

	err = tmpl.Execute(&buffer, spec)
	if err != nil {
		return "", "", "", err
	}
	return buffer.String(), reqCode, resCode, err
}
