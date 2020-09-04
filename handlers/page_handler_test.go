package handlers

import (
	"net.vikesh/goshop/config"
	"net/http/httptest"
	"testing"
)

var tcs = []config.TemplateConfiguration{
	{ParseOnce: false, TemplateDirectory: "./", Favicon: "./page_handler_test.go"},
	{ParseOnce: true, TemplateDirectory: "./", Favicon: "vxfdsfdsfewdsfsd"},
}

func Test_parse(t *testing.T) {

	tests := []struct {
		name           string
		args           string
		wantErr        bool
		templateConfig config.TemplateConfiguration
	}{
		{"should parse template", "test_template.gohtml", false, tcs[0]},
		{"should parse and cache template", "test_template.gohtml", false, tcs[1]},
		{"should fail template with single parse", "test_template2.gohtml", true, tcs[1]},
		{"should fail to parse template when caching is enabled", "test_template2.gohtml", true, tcs[0]},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc = tt.templateConfig
			tm, err := parse(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("parse() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil && tm == nil {
				t.Errorf("parse() wanted template got = %v", tm)
			}

			if tt.templateConfig.ParseOnce && err == nil && templates[tt.args] == nil {
				t.Errorf("template not put in map")
			}

			if !tt.templateConfig.ParseOnce && err == nil && templates[tt.args] != nil {
				t.Errorf("template put in map when configured to parse every time")
			}
		})
	}
}

func TestCreateHandlers(t *testing.T) {
	type args struct {
		c config.TemplateConfiguration
	}
	tests := []struct {
		name string
		args args
	}{
		{"template map should not be empty", args{tcs[0]}},
		{"template map should not be empty", args{tcs[1]}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateHandlers(tt.args.c); len(got) == 0 {
				t.Errorf("CreateHandlers() = %v, want %v", got, "not nil map")
			}
		})
	}
}

func Test_faviconHandler(t *testing.T) {

	tests := []struct {
		name     string
		template config.TemplateConfiguration
		url      string
		method   string
		status   int
	}{
		{"want 200", tcs[0], "/favicon.ico", "GET", 200},
		{"want 200", tcs[1], "/favicon.ico", "GET", 404},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc = tt.template
			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			req := httptest.NewRequest(tt.method, tt.url, nil)
			w := httptest.NewRecorder()
			faviconHandler(w, req)
			resp := w.Result()
			if resp.StatusCode != tt.status {
				t.Errorf(tt.name, "wanted %v, got %v", tt.status, tt.status)
			}
		})
	}
}
