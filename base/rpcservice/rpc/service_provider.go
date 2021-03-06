// Autogenerated by Thrift Compiler (0.9.1)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package rpc

import (
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"math"
)

// (needed to ensure safety because of naive import list construction.)
var _ = math.MinInt32
var _ = thrift.ZERO
var _ = fmt.Printf

type ServiceProvider interface {
	// Parameters:
	//  - Name
	//  - Input
	//  - Session
	//  - Timeout
	Request(name string, input string, session string, timeout int64) (r string, err error)
	// Parameters:
	//  - Name
	//  - Input
	//  - Data
	//  - Timeout
	Send(name string, input string, data []byte, timeout int64) (r string, err error)
	// Parameters:
	//  - Input
	//  - Timeout
	Heartbeat(input string, timeout int64) (r string, err error)
	// Parameters:
	//  - Name
	//  - Input
	//  - Timeout
	Get(name string, input string, timeout int64) (r []byte, err error)
}

type ServiceProviderClient struct {
	Transport       thrift.TTransport
	ProtocolFactory thrift.TProtocolFactory
	InputProtocol   thrift.TProtocol
	OutputProtocol  thrift.TProtocol
	SeqId           int32
}

func NewServiceProviderClientFactory(t thrift.TTransport, f thrift.TProtocolFactory) *ServiceProviderClient {
	return &ServiceProviderClient{Transport: t,
		ProtocolFactory: f,
		InputProtocol:   f.GetProtocol(t),
		OutputProtocol:  f.GetProtocol(t),
		SeqId:           0,
	}
}

func NewServiceProviderClientProtocol(t thrift.TTransport, iprot thrift.TProtocol, oprot thrift.TProtocol) *ServiceProviderClient {
	return &ServiceProviderClient{Transport: t,
		ProtocolFactory: nil,
		InputProtocol:   iprot,
		OutputProtocol:  oprot,
		SeqId:           0,
	}
}

// Parameters:
//  - Name
//  - Input
//  - Session
//  - Timeout
func (p *ServiceProviderClient) Request(name string, input string, session string, timeout int64) (r string, err error) {
	if err = p.sendRequest(name, input, session, timeout); err != nil {
		return
	}
	return p.recvRequest()
}

func (p *ServiceProviderClient) sendRequest(name string, input string, session string, timeout int64) (err error) {
	oprot := p.OutputProtocol
	if oprot == nil {
		oprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.OutputProtocol = oprot
	}
	p.SeqId++
	oprot.WriteMessageBegin("Request", thrift.CALL, p.SeqId)
	args0 := NewRequestArgs()
	args0.Name = name
	args0.Input = input
	args0.Session = session
	args0.Timeout = timeout
	err = args0.Write(oprot)
	oprot.WriteMessageEnd()
	oprot.Flush()
	return
}

func (p *ServiceProviderClient) recvRequest() (value string, err error) {
	iprot := p.InputProtocol
	if iprot == nil {
		iprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.InputProtocol = iprot
	}
	_, mTypeId, seqId, err := iprot.ReadMessageBegin()
	if err != nil {
		return
	}
	if mTypeId == thrift.EXCEPTION {
		error2 := thrift.NewTApplicationException(thrift.UNKNOWN_APPLICATION_EXCEPTION, "Unknown Exception")
		var error3 error
		error3, err = error2.Read(iprot)
		if err != nil {
			return
		}
		if err = iprot.ReadMessageEnd(); err != nil {
			return
		}
		err = error3
		return
	}
	if p.SeqId != seqId {
		err = thrift.NewTApplicationException(thrift.BAD_SEQUENCE_ID, "ping failed: out of sequence response")
		return
	}
	result1 := NewRequestResult()
	err = result1.Read(iprot)
	iprot.ReadMessageEnd()
	value = result1.Success
	return
}

// Parameters:
//  - Name
//  - Input
//  - Data
//  - Timeout
func (p *ServiceProviderClient) Send(name string, input string, data []byte, timeout int64) (r string, err error) {
	if err = p.sendSend(name, input, data, timeout); err != nil {
		return
	}
	return p.recvSend()
}

func (p *ServiceProviderClient) sendSend(name string, input string, data []byte, timeout int64) (err error) {
	oprot := p.OutputProtocol
	if oprot == nil {
		oprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.OutputProtocol = oprot
	}
	p.SeqId++
	oprot.WriteMessageBegin("Send", thrift.CALL, p.SeqId)
	args4 := NewSendArgs()
	args4.Name = name
	args4.Input = input
	args4.Data = data
	args4.Timeout = timeout
	err = args4.Write(oprot)
	oprot.WriteMessageEnd()
	oprot.Flush()
	return
}

func (p *ServiceProviderClient) recvSend() (value string, err error) {
	iprot := p.InputProtocol
	if iprot == nil {
		iprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.InputProtocol = iprot
	}
	_, mTypeId, seqId, err := iprot.ReadMessageBegin()
	if err != nil {
		return
	}
	if mTypeId == thrift.EXCEPTION {
		error6 := thrift.NewTApplicationException(thrift.UNKNOWN_APPLICATION_EXCEPTION, "Unknown Exception")
		var error7 error
		error7, err = error6.Read(iprot)
		if err != nil {
			return
		}
		if err = iprot.ReadMessageEnd(); err != nil {
			return
		}
		err = error7
		return
	}
	if p.SeqId != seqId {
		err = thrift.NewTApplicationException(thrift.BAD_SEQUENCE_ID, "ping failed: out of sequence response")
		return
	}
	result5 := NewSendResult()
	err = result5.Read(iprot)
	iprot.ReadMessageEnd()
	value = result5.Success
	return
}

// Parameters:
//  - Input
//  - Timeout
func (p *ServiceProviderClient) Heartbeat(input string, timeout int64) (r string, err error) {
	if err = p.sendHeartbeat(input, timeout); err != nil {
		return
	}
	return p.recvHeartbeat()
}

func (p *ServiceProviderClient) sendHeartbeat(input string, timeout int64) (err error) {
	oprot := p.OutputProtocol
	if oprot == nil {
		oprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.OutputProtocol = oprot
	}
	p.SeqId++
	oprot.WriteMessageBegin("Heartbeat", thrift.CALL, p.SeqId)
	args8 := NewHeartbeatArgs()
	args8.Input = input
	args8.Timeout = timeout
	err = args8.Write(oprot)
	oprot.WriteMessageEnd()
	oprot.Flush()
	return
}

func (p *ServiceProviderClient) recvHeartbeat() (value string, err error) {
	iprot := p.InputProtocol
	if iprot == nil {
		iprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.InputProtocol = iprot
	}
	_, mTypeId, seqId, err := iprot.ReadMessageBegin()
	if err != nil {
		return
	}
	if mTypeId == thrift.EXCEPTION {
		error10 := thrift.NewTApplicationException(thrift.UNKNOWN_APPLICATION_EXCEPTION, "Unknown Exception")
		var error11 error
		error11, err = error10.Read(iprot)
		if err != nil {
			return
		}
		if err = iprot.ReadMessageEnd(); err != nil {
			return
		}
		err = error11
		return
	}
	if p.SeqId != seqId {
		err = thrift.NewTApplicationException(thrift.BAD_SEQUENCE_ID, "ping failed: out of sequence response")
		return
	}
	result9 := NewHeartbeatResult()
	err = result9.Read(iprot)
	iprot.ReadMessageEnd()
	value = result9.Success
	return
}

// Parameters:
//  - Name
//  - Input
//  - Timeout
func (p *ServiceProviderClient) Get(name string, input string, timeout int64) (r []byte, err error) {
	if err = p.sendGet(name, input, timeout); err != nil {
		return
	}
	return p.recvGet()
}

func (p *ServiceProviderClient) sendGet(name string, input string, timeout int64) (err error) {
	oprot := p.OutputProtocol
	if oprot == nil {
		oprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.OutputProtocol = oprot
	}
	p.SeqId++
	oprot.WriteMessageBegin("Get", thrift.CALL, p.SeqId)
	args12 := NewGetArgs()
	args12.Name = name
	args12.Input = input
	args12.Timeout = timeout
	err = args12.Write(oprot)
	oprot.WriteMessageEnd()
	oprot.Flush()
	return
}

func (p *ServiceProviderClient) recvGet() (value []byte, err error) {
	iprot := p.InputProtocol
	if iprot == nil {
		iprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.InputProtocol = iprot
	}
	_, mTypeId, seqId, err := iprot.ReadMessageBegin()
	if err != nil {
		return
	}
	if mTypeId == thrift.EXCEPTION {
		error14 := thrift.NewTApplicationException(thrift.UNKNOWN_APPLICATION_EXCEPTION, "Unknown Exception")
		var error15 error
		error15, err = error14.Read(iprot)
		if err != nil {
			return
		}
		if err = iprot.ReadMessageEnd(); err != nil {
			return
		}
		err = error15
		return
	}
	if p.SeqId != seqId {
		err = thrift.NewTApplicationException(thrift.BAD_SEQUENCE_ID, "ping failed: out of sequence response")
		return
	}
	result13 := NewGetResult()
	err = result13.Read(iprot)
	iprot.ReadMessageEnd()
	value = result13.Success
	return
}

type ServiceProviderProcessor struct {
	processorMap map[string]thrift.TProcessorFunction
	handler      ServiceProvider
}

func (p *ServiceProviderProcessor) AddToProcessorMap(key string, processor thrift.TProcessorFunction) {
	p.processorMap[key] = processor
}

func (p *ServiceProviderProcessor) GetProcessorFunction(key string) (processor thrift.TProcessorFunction, ok bool) {
	processor, ok = p.processorMap[key]
	return processor, ok
}

func (p *ServiceProviderProcessor) ProcessorMap() map[string]thrift.TProcessorFunction {
	return p.processorMap
}

func NewServiceProviderProcessor(handler ServiceProvider) *ServiceProviderProcessor {

	self16 := &ServiceProviderProcessor{handler: handler, processorMap: make(map[string]thrift.TProcessorFunction)}
	self16.processorMap["Request"] = &serviceProviderProcessorRequest{handler: handler}
	self16.processorMap["Send"] = &serviceProviderProcessorSend{handler: handler}
	self16.processorMap["Heartbeat"] = &serviceProviderProcessorHeartbeat{handler: handler}
	self16.processorMap["Get"] = &serviceProviderProcessorGet{handler: handler}
	return self16
}

func (p *ServiceProviderProcessor) Process(iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	name, _, seqId, err := iprot.ReadMessageBegin()
	if err != nil {
		return false, err
	}
	if processor, ok := p.GetProcessorFunction(name); ok {
		return processor.Process(seqId, iprot, oprot)
	}
	iprot.Skip(thrift.STRUCT)
	iprot.ReadMessageEnd()
	x17 := thrift.NewTApplicationException(thrift.UNKNOWN_METHOD, "Unknown function "+name)
	oprot.WriteMessageBegin(name, thrift.EXCEPTION, seqId)
	x17.Write(oprot)
	oprot.WriteMessageEnd()
	oprot.Flush()
	return false, x17

}

type serviceProviderProcessorRequest struct {
	handler ServiceProvider
}

func (p *serviceProviderProcessorRequest) Process(seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	args := NewRequestArgs()
	if err = args.Read(iprot); err != nil {
		iprot.ReadMessageEnd()
		x := thrift.NewTApplicationException(thrift.PROTOCOL_ERROR, err.Error())
		oprot.WriteMessageBegin("Request", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush()
		return
	}
	iprot.ReadMessageEnd()
	result := NewRequestResult()
	if result.Success, err = p.handler.Request(args.Name, args.Input, args.Session, args.Timeout); err != nil {
		x := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, "Internal error processing Request: "+err.Error())
		oprot.WriteMessageBegin("Request", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush()
		return
	}
	if err2 := oprot.WriteMessageBegin("Request", thrift.REPLY, seqId); err2 != nil {
		err = err2
	}
	if err2 := result.Write(oprot); err == nil && err2 != nil {
		err = err2
	}
	if err2 := oprot.WriteMessageEnd(); err == nil && err2 != nil {
		err = err2
	}
	if err2 := oprot.Flush(); err == nil && err2 != nil {
		err = err2
	}
	if err != nil {
		return
	}
	return true, err
}

type serviceProviderProcessorSend struct {
	handler ServiceProvider
}

func (p *serviceProviderProcessorSend) Process(seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	args := NewSendArgs()
	if err = args.Read(iprot); err != nil {
		iprot.ReadMessageEnd()
		x := thrift.NewTApplicationException(thrift.PROTOCOL_ERROR, err.Error())
		oprot.WriteMessageBegin("Send", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush()
		return
	}
	iprot.ReadMessageEnd()
	result := NewSendResult()
	if result.Success, err = p.handler.Send(args.Name, args.Input, args.Data, args.Timeout); err != nil {
		x := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, "Internal error processing Send: "+err.Error())
		oprot.WriteMessageBegin("Send", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush()
		return
	}
	if err2 := oprot.WriteMessageBegin("Send", thrift.REPLY, seqId); err2 != nil {
		err = err2
	}
	if err2 := result.Write(oprot); err == nil && err2 != nil {
		err = err2
	}
	if err2 := oprot.WriteMessageEnd(); err == nil && err2 != nil {
		err = err2
	}
	if err2 := oprot.Flush(); err == nil && err2 != nil {
		err = err2
	}
	if err != nil {
		return
	}
	return true, err
}

type serviceProviderProcessorHeartbeat struct {
	handler ServiceProvider
}

func (p *serviceProviderProcessorHeartbeat) Process(seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	args := NewHeartbeatArgs()
	if err = args.Read(iprot); err != nil {
		iprot.ReadMessageEnd()
		x := thrift.NewTApplicationException(thrift.PROTOCOL_ERROR, err.Error())
		oprot.WriteMessageBegin("Heartbeat", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush()
		return
	}
	iprot.ReadMessageEnd()
	result := NewHeartbeatResult()
	if result.Success, err = p.handler.Heartbeat(args.Input, args.Timeout); err != nil {
		x := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, "Internal error processing Heartbeat: "+err.Error())
		oprot.WriteMessageBegin("Heartbeat", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush()
		return
	}
	if err2 := oprot.WriteMessageBegin("Heartbeat", thrift.REPLY, seqId); err2 != nil {
		err = err2
	}
	if err2 := result.Write(oprot); err == nil && err2 != nil {
		err = err2
	}
	if err2 := oprot.WriteMessageEnd(); err == nil && err2 != nil {
		err = err2
	}
	if err2 := oprot.Flush(); err == nil && err2 != nil {
		err = err2
	}
	if err != nil {
		return
	}
	return true, err
}

type serviceProviderProcessorGet struct {
	handler ServiceProvider
}

func (p *serviceProviderProcessorGet) Process(seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	args := NewGetArgs()
	if err = args.Read(iprot); err != nil {
		iprot.ReadMessageEnd()
		x := thrift.NewTApplicationException(thrift.PROTOCOL_ERROR, err.Error())
		oprot.WriteMessageBegin("Get", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush()
		return
	}
	iprot.ReadMessageEnd()
	result := NewGetResult()
	if result.Success, err = p.handler.Get(args.Name, args.Input, args.Timeout); err != nil {
		x := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, "Internal error processing Get: "+err.Error())
		oprot.WriteMessageBegin("Get", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush()
		return
	}
	if err2 := oprot.WriteMessageBegin("Get", thrift.REPLY, seqId); err2 != nil {
		err = err2
	}
	if err2 := result.Write(oprot); err == nil && err2 != nil {
		err = err2
	}
	if err2 := oprot.WriteMessageEnd(); err == nil && err2 != nil {
		err = err2
	}
	if err2 := oprot.Flush(); err == nil && err2 != nil {
		err = err2
	}
	if err != nil {
		return
	}
	return true, err
}

// HELPER FUNCTIONS AND STRUCTURES

type RequestArgs struct {
	Name    string `thrift:"name,1"`
	Input   string `thrift:"input,2"`
	Session string `thrift:"session,3"`
	Timeout int64  `thrift:"timeout,4"`
}

func NewRequestArgs() *RequestArgs {
	return &RequestArgs{}
}

func (p *RequestArgs) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return fmt.Errorf("%T read error", p)
	}
	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return fmt.Errorf("%T field %d read error: %s", p, fieldId, err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.readField1(iprot); err != nil {
				return err
			}
		case 2:
			if err := p.readField2(iprot); err != nil {
				return err
			}
		case 3:
			if err := p.readField3(iprot); err != nil {
				return err
			}
		case 4:
			if err := p.readField4(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return fmt.Errorf("%T read struct end error: %s", p, err)
	}
	return nil
}

func (p *RequestArgs) readField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 1: %s")
	} else {
		p.Name = v
	}
	return nil
}

func (p *RequestArgs) readField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 2: %s")
	} else {
		p.Input = v
	}
	return nil
}

func (p *RequestArgs) readField3(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 3: %s")
	} else {
		p.Session = v
	}
	return nil
}

func (p *RequestArgs) readField4(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(); err != nil {
		return fmt.Errorf("error reading field 4: %s")
	} else {
		p.Timeout = v
	}
	return nil
}

func (p *RequestArgs) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("Request_args"); err != nil {
		return fmt.Errorf("%T write struct begin error: %s", p, err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := p.writeField2(oprot); err != nil {
		return err
	}
	if err := p.writeField3(oprot); err != nil {
		return err
	}
	if err := p.writeField4(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return fmt.Errorf("%T write field stop error: %s", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return fmt.Errorf("%T write struct stop error: %s", err)
	}
	return nil
}

func (p *RequestArgs) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("name", thrift.STRING, 1); err != nil {
		return fmt.Errorf("%T write field begin error 1:name: %s", p, err)
	}
	if err := oprot.WriteString(string(p.Name)); err != nil {
		return fmt.Errorf("%T.name (1) field write error: %s", p)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 1:name: %s", p, err)
	}
	return err
}

func (p *RequestArgs) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("input", thrift.STRING, 2); err != nil {
		return fmt.Errorf("%T write field begin error 2:input: %s", p, err)
	}
	if err := oprot.WriteString(string(p.Input)); err != nil {
		return fmt.Errorf("%T.input (2) field write error: %s", p)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 2:input: %s", p, err)
	}
	return err
}

func (p *RequestArgs) writeField3(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("session", thrift.STRING, 3); err != nil {
		return fmt.Errorf("%T write field begin error 3:session: %s", p, err)
	}
	if err := oprot.WriteString(string(p.Session)); err != nil {
		return fmt.Errorf("%T.session (3) field write error: %s", p)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 3:session: %s", p, err)
	}
	return err
}

func (p *RequestArgs) writeField4(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("timeout", thrift.I64, 4); err != nil {
		return fmt.Errorf("%T write field begin error 4:timeout: %s", p, err)
	}
	if err := oprot.WriteI64(int64(p.Timeout)); err != nil {
		return fmt.Errorf("%T.timeout (4) field write error: %s", p)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 4:timeout: %s", p, err)
	}
	return err
}

func (p *RequestArgs) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("RequestArgs(%+v)", *p)
}

type RequestResult struct {
	Success string `thrift:"success,0"`
}

func NewRequestResult() *RequestResult {
	return &RequestResult{}
}

func (p *RequestResult) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return fmt.Errorf("%T read error", p)
	}
	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return fmt.Errorf("%T field %d read error: %s", p, fieldId, err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 0:
			if err := p.readField0(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return fmt.Errorf("%T read struct end error: %s", p, err)
	}
	return nil
}

func (p *RequestResult) readField0(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 0: %s")
	} else {
		p.Success = v
	}
	return nil
}

func (p *RequestResult) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("Request_result"); err != nil {
		return fmt.Errorf("%T write struct begin error: %s", p, err)
	}
	switch {
	default:
		if err := p.writeField0(oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return fmt.Errorf("%T write field stop error: %s", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return fmt.Errorf("%T write struct stop error: %s", err)
	}
	return nil
}

func (p *RequestResult) writeField0(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("success", thrift.STRING, 0); err != nil {
		return fmt.Errorf("%T write field begin error 0:success: %s", p, err)
	}
	if err := oprot.WriteString(string(p.Success)); err != nil {
		return fmt.Errorf("%T.success (0) field write error: %s", p)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 0:success: %s", p, err)
	}
	return err
}

func (p *RequestResult) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("RequestResult(%+v)", *p)
}

type SendArgs struct {
	Name    string `thrift:"name,1"`
	Input   string `thrift:"input,2"`
	Data    []byte `thrift:"data,3"`
	Timeout int64  `thrift:"timeout,4"`
}

func NewSendArgs() *SendArgs {
	return &SendArgs{}
}

func (p *SendArgs) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return fmt.Errorf("%T read error", p)
	}
	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return fmt.Errorf("%T field %d read error: %s", p, fieldId, err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.readField1(iprot); err != nil {
				return err
			}
		case 2:
			if err := p.readField2(iprot); err != nil {
				return err
			}
		case 3:
			if err := p.readField3(iprot); err != nil {
				return err
			}
		case 4:
			if err := p.readField4(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return fmt.Errorf("%T read struct end error: %s", p, err)
	}
	return nil
}

func (p *SendArgs) readField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 1: %s")
	} else {
		p.Name = v
	}
	return nil
}

func (p *SendArgs) readField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 2: %s")
	} else {
		p.Input = v
	}
	return nil
}

func (p *SendArgs) readField3(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadBinary(); err != nil {
		return fmt.Errorf("error reading field 3: %s")
	} else {
		p.Data = v
	}
	return nil
}

func (p *SendArgs) readField4(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(); err != nil {
		return fmt.Errorf("error reading field 4: %s")
	} else {
		p.Timeout = v
	}
	return nil
}

func (p *SendArgs) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("Send_args"); err != nil {
		return fmt.Errorf("%T write struct begin error: %s", p, err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := p.writeField2(oprot); err != nil {
		return err
	}
	if err := p.writeField3(oprot); err != nil {
		return err
	}
	if err := p.writeField4(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return fmt.Errorf("%T write field stop error: %s", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return fmt.Errorf("%T write struct stop error: %s", err)
	}
	return nil
}

func (p *SendArgs) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("name", thrift.STRING, 1); err != nil {
		return fmt.Errorf("%T write field begin error 1:name: %s", p, err)
	}
	if err := oprot.WriteString(string(p.Name)); err != nil {
		return fmt.Errorf("%T.name (1) field write error: %s", p)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 1:name: %s", p, err)
	}
	return err
}

func (p *SendArgs) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("input", thrift.STRING, 2); err != nil {
		return fmt.Errorf("%T write field begin error 2:input: %s", p, err)
	}
	if err := oprot.WriteString(string(p.Input)); err != nil {
		return fmt.Errorf("%T.input (2) field write error: %s", p)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 2:input: %s", p, err)
	}
	return err
}

func (p *SendArgs) writeField3(oprot thrift.TProtocol) (err error) {
	if p.Data != nil {
		if err := oprot.WriteFieldBegin("data", thrift.BINARY, 3); err != nil {
			return fmt.Errorf("%T write field begin error 3:data: %s", p, err)
		}
		if err := oprot.WriteBinary(p.Data); err != nil {
			return fmt.Errorf("%T.data (3) field write error: %s", p)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return fmt.Errorf("%T write field end error 3:data: %s", p, err)
		}
	}
	return err
}

func (p *SendArgs) writeField4(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("timeout", thrift.I64, 4); err != nil {
		return fmt.Errorf("%T write field begin error 4:timeout: %s", p, err)
	}
	if err := oprot.WriteI64(int64(p.Timeout)); err != nil {
		return fmt.Errorf("%T.timeout (4) field write error: %s", p)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 4:timeout: %s", p, err)
	}
	return err
}

func (p *SendArgs) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("SendArgs(%+v)", *p)
}

type SendResult struct {
	Success string `thrift:"success,0"`
}

func NewSendResult() *SendResult {
	return &SendResult{}
}

func (p *SendResult) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return fmt.Errorf("%T read error", p)
	}
	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return fmt.Errorf("%T field %d read error: %s", p, fieldId, err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 0:
			if err := p.readField0(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return fmt.Errorf("%T read struct end error: %s", p, err)
	}
	return nil
}

func (p *SendResult) readField0(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 0: %s")
	} else {
		p.Success = v
	}
	return nil
}

func (p *SendResult) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("Send_result"); err != nil {
		return fmt.Errorf("%T write struct begin error: %s", p, err)
	}
	switch {
	default:
		if err := p.writeField0(oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return fmt.Errorf("%T write field stop error: %s", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return fmt.Errorf("%T write struct stop error: %s", err)
	}
	return nil
}

func (p *SendResult) writeField0(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("success", thrift.STRING, 0); err != nil {
		return fmt.Errorf("%T write field begin error 0:success: %s", p, err)
	}
	if err := oprot.WriteString(string(p.Success)); err != nil {
		return fmt.Errorf("%T.success (0) field write error: %s", p)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 0:success: %s", p, err)
	}
	return err
}

func (p *SendResult) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("SendResult(%+v)", *p)
}

type HeartbeatArgs struct {
	Input   string `thrift:"input,1"`
	Timeout int64  `thrift:"timeout,2"`
}

func NewHeartbeatArgs() *HeartbeatArgs {
	return &HeartbeatArgs{}
}

func (p *HeartbeatArgs) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return fmt.Errorf("%T read error", p)
	}
	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return fmt.Errorf("%T field %d read error: %s", p, fieldId, err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.readField1(iprot); err != nil {
				return err
			}
		case 2:
			if err := p.readField2(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return fmt.Errorf("%T read struct end error: %s", p, err)
	}
	return nil
}

func (p *HeartbeatArgs) readField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 1: %s")
	} else {
		p.Input = v
	}
	return nil
}

func (p *HeartbeatArgs) readField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(); err != nil {
		return fmt.Errorf("error reading field 2: %s")
	} else {
		p.Timeout = v
	}
	return nil
}

func (p *HeartbeatArgs) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("Heartbeat_args"); err != nil {
		return fmt.Errorf("%T write struct begin error: %s", p, err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := p.writeField2(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return fmt.Errorf("%T write field stop error: %s", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return fmt.Errorf("%T write struct stop error: %s", err)
	}
	return nil
}

func (p *HeartbeatArgs) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("input", thrift.STRING, 1); err != nil {
		return fmt.Errorf("%T write field begin error 1:input: %s", p, err)
	}
	if err := oprot.WriteString(string(p.Input)); err != nil {
		return fmt.Errorf("%T.input (1) field write error: %s", p)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 1:input: %s", p, err)
	}
	return err
}

func (p *HeartbeatArgs) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("timeout", thrift.I64, 2); err != nil {
		return fmt.Errorf("%T write field begin error 2:timeout: %s", p, err)
	}
	if err := oprot.WriteI64(int64(p.Timeout)); err != nil {
		return fmt.Errorf("%T.timeout (2) field write error: %s", p)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 2:timeout: %s", p, err)
	}
	return err
}

func (p *HeartbeatArgs) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("HeartbeatArgs(%+v)", *p)
}

type HeartbeatResult struct {
	Success string `thrift:"success,0"`
}

func NewHeartbeatResult() *HeartbeatResult {
	return &HeartbeatResult{}
}

func (p *HeartbeatResult) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return fmt.Errorf("%T read error", p)
	}
	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return fmt.Errorf("%T field %d read error: %s", p, fieldId, err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 0:
			if err := p.readField0(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return fmt.Errorf("%T read struct end error: %s", p, err)
	}
	return nil
}

func (p *HeartbeatResult) readField0(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 0: %s")
	} else {
		p.Success = v
	}
	return nil
}

func (p *HeartbeatResult) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("Heartbeat_result"); err != nil {
		return fmt.Errorf("%T write struct begin error: %s", p, err)
	}
	switch {
	default:
		if err := p.writeField0(oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return fmt.Errorf("%T write field stop error: %s", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return fmt.Errorf("%T write struct stop error: %s", err)
	}
	return nil
}

func (p *HeartbeatResult) writeField0(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("success", thrift.STRING, 0); err != nil {
		return fmt.Errorf("%T write field begin error 0:success: %s", p, err)
	}
	if err := oprot.WriteString(string(p.Success)); err != nil {
		return fmt.Errorf("%T.success (0) field write error: %s", p)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 0:success: %s", p, err)
	}
	return err
}

func (p *HeartbeatResult) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("HeartbeatResult(%+v)", *p)
}

type GetArgs struct {
	Name    string `thrift:"name,1"`
	Input   string `thrift:"input,2"`
	Timeout int64  `thrift:"timeout,3"`
}

func NewGetArgs() *GetArgs {
	return &GetArgs{}
}

func (p *GetArgs) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return fmt.Errorf("%T read error", p)
	}
	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return fmt.Errorf("%T field %d read error: %s", p, fieldId, err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.readField1(iprot); err != nil {
				return err
			}
		case 2:
			if err := p.readField2(iprot); err != nil {
				return err
			}
		case 3:
			if err := p.readField3(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return fmt.Errorf("%T read struct end error: %s", p, err)
	}
	return nil
}

func (p *GetArgs) readField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 1: %s")
	} else {
		p.Name = v
	}
	return nil
}

func (p *GetArgs) readField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 2: %s")
	} else {
		p.Input = v
	}
	return nil
}

func (p *GetArgs) readField3(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(); err != nil {
		return fmt.Errorf("error reading field 3: %s")
	} else {
		p.Timeout = v
	}
	return nil
}

func (p *GetArgs) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("Get_args"); err != nil {
		return fmt.Errorf("%T write struct begin error: %s", p, err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := p.writeField2(oprot); err != nil {
		return err
	}
	if err := p.writeField3(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return fmt.Errorf("%T write field stop error: %s", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return fmt.Errorf("%T write struct stop error: %s", err)
	}
	return nil
}

func (p *GetArgs) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("name", thrift.STRING, 1); err != nil {
		return fmt.Errorf("%T write field begin error 1:name: %s", p, err)
	}
	if err := oprot.WriteString(string(p.Name)); err != nil {
		return fmt.Errorf("%T.name (1) field write error: %s", p)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 1:name: %s", p, err)
	}
	return err
}

func (p *GetArgs) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("input", thrift.STRING, 2); err != nil {
		return fmt.Errorf("%T write field begin error 2:input: %s", p, err)
	}
	if err := oprot.WriteString(string(p.Input)); err != nil {
		return fmt.Errorf("%T.input (2) field write error: %s", p)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 2:input: %s", p, err)
	}
	return err
}

func (p *GetArgs) writeField3(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("timeout", thrift.I64, 3); err != nil {
		return fmt.Errorf("%T write field begin error 3:timeout: %s", p, err)
	}
	if err := oprot.WriteI64(int64(p.Timeout)); err != nil {
		return fmt.Errorf("%T.timeout (3) field write error: %s", p)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 3:timeout: %s", p, err)
	}
	return err
}

func (p *GetArgs) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("GetArgs(%+v)", *p)
}

type GetResult struct {
	Success []byte `thrift:"success,0"`
}

func NewGetResult() *GetResult {
	return &GetResult{}
}

func (p *GetResult) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return fmt.Errorf("%T read error", p)
	}
	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return fmt.Errorf("%T field %d read error: %s", p, fieldId, err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 0:
			if err := p.readField0(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return fmt.Errorf("%T read struct end error: %s", p, err)
	}
	return nil
}

func (p *GetResult) readField0(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadBinary(); err != nil {
		return fmt.Errorf("error reading field 0: %s")
	} else {
		p.Success = v
	}
	return nil
}

func (p *GetResult) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("Get_result"); err != nil {
		return fmt.Errorf("%T write struct begin error: %s", p, err)
	}
	switch {
	default:
		if err := p.writeField0(oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return fmt.Errorf("%T write field stop error: %s", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return fmt.Errorf("%T write struct stop error: %s", err)
	}
	return nil
}

func (p *GetResult) writeField0(oprot thrift.TProtocol) (err error) {
	if p.Success != nil {
		if err := oprot.WriteFieldBegin("success", thrift.BINARY, 0); err != nil {
			return fmt.Errorf("%T write field begin error 0:success: %s", p, err)
		}
		if err := oprot.WriteBinary(p.Success); err != nil {
			return fmt.Errorf("%T.success (0) field write error: %s", p)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return fmt.Errorf("%T write field end error 0:success: %s", p, err)
		}
	}
	return err
}

func (p *GetResult) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("GetResult(%+v)", *p)
}
