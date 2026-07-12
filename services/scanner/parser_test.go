package scanner

import "testing"

func TestParseMerge(t *testing.T) {
	input := []byte(`[
		{"ip":"10.0.0.5","timestamp":"1","ports":[{"port":22,"proto":"tcp","status":"open","reason":"syn-ack","ttl":64}]},
		{"ip":"10.0.0.5","timestamp":"2","ports":[{"port":22,"proto":"tcp","service":{"name":"ssh","banner":"SSH-2.0-OpenSSH_8.9"}}]},
		{"ip":"10.0.0.5","timestamp":"3","ports":[{"port":80,"proto":"tcp","status":"open","reason":"syn-ack","ttl":64}]},
		{"finished":1}
	]`)

	ports, err := parseMasscanJSON(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(ports) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(ports))
	}

	ssh := ports[0]
	if ssh.Port != 22 || ssh.Service != "ssh" || ssh.Banner != "SSH-2.0-OpenSSH_8.9" {
		t.Errorf("port 22 merge failed: %+v", ssh)
	}

	http := ports[1]
	if http.Port != 80 || http.Service != "" {
		t.Errorf("port 80 unexpected: %+v", http)
	}
}

func TestParseEmpty(t *testing.T) {
	ports, err := parseMasscanJSON([]byte("[]"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ports) != 0 {
		t.Errorf("expected 0 ports, got %d", len(ports))
	}
}
