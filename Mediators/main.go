package main

import (
        "bytes"
        "encoding/json"
        "errors"
        "fmt"
        "io"
        "log"
        "sync"

        "github.com/oliveagle/jsonpath"
        "github.com/lestrrat-go/libxml2"
)

// Envelope represents the message envelope.
type Envelope struct {
        Headers     map[string][]string
        Properties  map[string]interface{}
        Payload     Payload
        Attachments map[string]Payload
        Error       error
        mu          sync.Mutex // For concurrency safety
}

// Payload is an interface that represents the message payload.
type Payload interface {
        Type() string
        Reader() io.Reader
        Bytes() ([]byte, error)
        Schema() interface{}
        SetSchema(interface{})
        SetReader(io.Reader)
        SetBytes([]byte)
}

// DeferredPayload stores the raw data and parses it on demand.
type DeferredPayload struct {
        dataType string
        data     interface{} // io.Reader or []byte
        schema   interface{}
        bytes    []byte
        mu       sync.Mutex
}

func (p *DeferredPayload) Type() string {
        return p.dataType
}

func (p *DeferredPayload) Reader() io.Reader {
        p.mu.Lock()
        defer p.mu.Unlock()
        if r, ok := p.data.(io.Reader); ok {
                return r
        }
        if b, ok := p.data.([]byte); ok{
                return bytes.NewReader(b)
        }
        return nil
}

func (p *DeferredPayload) Bytes() ([]byte, error) {
        p.mu.Lock()
        defer p.mu.Unlock()

        if p.bytes != nil {
                return p.bytes, nil
        }

        if r, ok := p.data.(io.Reader); ok {
                b, err := io.ReadAll(r)
                if err != nil {
                        return nil, err
                }
                p.bytes = b
                return b, nil
        }

        if b, ok := p.data.([]byte); ok{
                p.bytes = b;
                return b, nil;
        }

        return nil, errors.New("no payload data")
}

func (p *DeferredPayload) Schema() interface{} {
        return p.schema
}

func (p *DeferredPayload) SetSchema(schema interface{}) {
        p.schema = schema
}

func (p *DeferredPayload) SetReader(r io.Reader) {
        p.mu.Lock()
        defer p.mu.Unlock()
        p.data = r
        p.bytes = nil
}

func (p *DeferredPayload) SetBytes(b []byte) {
        p.mu.Lock()
        defer p.mu.Unlock()
        p.data = b
        p.bytes = b
}

// MemoryPayload stores the parsed data in memory.
type MemoryPayload struct {
        dataType string
        data     interface{}
        schema   interface{}
}

func (p *MemoryPayload) Type() string {
        return p.dataType
}

func (p *MemoryPayload) Reader() io.Reader {
        b, err := json.Marshal(p.data)
        if err != nil {
                return nil
        }
        return bytes.NewReader(b)
}

func (p *MemoryPayload) Bytes() ([]byte, error) {
        return json.Marshal(p.data)
}

func (p *MemoryPayload) Schema() interface{} {
        return p.schema
}

func (p *MemoryPayload) SetSchema(schema interface{}) {
        p.schema = schema
}

func (p *MemoryPayload) SetReader(r io.Reader) {
        //Not implemented for memory payload
}

func (p *MemoryPayload) SetBytes(b []byte) {
        //Not implemented for memory payload
}

// JSONPathMediator applies a JSONPath expression to the payload.
func JSONPathMediator(env *Envelope, jsonPathExpr string) error {
        payloadBytes, err := env.Payload.Bytes()
        if err != nil {
                return err
        }

        var jsonData interface{}
        if err := json.Unmarshal(payloadBytes, &jsonData); err != nil {
                return err
        }

        result, err := jsonpath.JsonPathLookup(jsonData, jsonPathExpr)
        if err != nil {
                return err
        }

        env.Payload = &MemoryPayload{dataType: "application/json", data: result}
        return nil
}

// XPathMediator applies an XPath expression to an XML payload.
func XPathMediator(env *Envelope, xpathExpr string) error {
        payloadBytes, err := env.Payload.Bytes()
        if err != nil {
                return err
        }

        doc, err := libxml2.Parse(payloadBytes)
        if err != nil {
                return err
        }
        defer doc.Free()

        result, err := doc.Find(xpathExpr)
        if err != nil {
                return err
        }

        nodes := result.NodeList()
        if len(nodes) > 0 {
                env.Payload = &DeferredPayload{dataType: "application/xml", data: []byte(nodes[0].String())} // Simplified.
                return nil
        }

        return nil;
}

// PropertyMediator sets or gets a property in the envelope.
func PropertyMediator(env *Envelope, propertyName string, propertyValue interface{}, set bool) error {
        if env.Properties == nil {
                env.Properties = make(map[string]interface{})
        }

        if set {
                env.Properties[propertyName] = propertyValue
        } else {
                if value, ok := env.Properties[propertyName]; ok {
                        env.Payload = &MemoryPayload{dataType: "application/json", data: value}
                } else {
                        return errors.New("property not found")
                }
        }
        return nil
}

func main() {
        jsonPayload := []byte(`{"name": "John", "age": 30, "address": {"street": "123 Main St"}}`)
        xmlPayload := []byte(`<root><element>data</element></root>`)

        env := &Envelope{
                Payload: &DeferredPayload{dataType: "application/json", data: jsonPayload},
                Headers: make(map[string][]string),
                Properties: make(map[string]interface{}),
        }

        // Sequence of mediators
        if err := JSONPathMediator(env, "$.address.street"); err != nil {
                log.Fatalf("JSONPathMediator error: %v", err)
        }

        env.Payload = &DeferredPayload{dataType: "application/xml", data: xmlPayload}

        if err := XPathMediator(env, "/root/element/text()"); err != nil {
                log.Fatalf("XPathMediator error: %v", err)
        }

        if err := PropertyMediator(env, "testProperty", "testValue", true); err != nil {
                log.Fatalf("PropertyMediator set error: %v", err)
        }

        if err := PropertyMediator(env, "testProperty", nil, false); err != nil {
                log.Fatalf("PropertyMediator get error: %v", err)
        }

        finalPayload, err := env.Payload.Bytes()
        if err != nil {
                log.Fatalf("Final payload error: %v", err)
        }

        fmt.Println(string(finalPayload))
}