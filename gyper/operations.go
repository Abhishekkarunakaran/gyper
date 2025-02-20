package gyper

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
)

func (c *Context) Bind(dataField any) error {
	switch c.Request.Header["Content-Type"] {
	case "application/json":	
		if err := json.NewDecoder(c.Request.Body).Decode(dataField); err != nil {
			return fmt.Errorf("failed to bind the request body to the provided struct")
		}
	case "application/xml":
		if err := xml.NewDecoder(c.Request.Body).Decode(dataField); err != nil {
			return fmt.Errorf("failed to bind the request body to the provided struct")
		}
	
	}	

return nil

}

func (c *Context) JSON(statusCode int, data any) error {

	response,err := json.Marshal(data)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	statusString := http.StatusText(statusCode)
	switch c.Request.Protocol {
	case HTTP1:
		c.conn.Write([]byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n",statusCode, statusString)))
	case HTTP2:
		c.conn.Write([]byte(fmt.Sprintf("HTTP/2 %d %s\r\n",statusCode, statusString)))
	}
	c.conn.Write([]byte("Content-Length: " + fmt.Sprint(len(response)) + "\r\n"))
	c.conn.Write([]byte("Content-Type: application/json\r\n\r\n"))
	c.conn.Write(response)

	return nil
}

func (c *Context) XML(statusCode int, data any) error {

	response,err := xml.MarshalIndent(data, " ","	")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	statusString := http.StatusText(statusCode)
	switch c.Request.Protocol {
	case HTTP1:
		c.conn.Write([]byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n",statusCode, statusString)))
	case HTTP2:
		c.conn.Write([]byte(fmt.Sprintf("HTTP/2 %d %s\r\n",statusCode, statusString)))
	}
	c.conn.Write([]byte("Content-Length: " + fmt.Sprint(len(response)) + "\r\n"))
	c.conn.Write([]byte("Content-Type: application/xml\r\n\r\n"))
	c.conn.Write(response)

	return nil
}
