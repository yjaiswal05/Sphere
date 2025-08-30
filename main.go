package main

import (
        "encoding/json"
        "fmt"
        "html/template"
        "log"
        "net"
        "net/http"
        "os"
        "path/filepath"
        "strings"
)

// VerificationResult represents the result of domain verification
type VerificationResult struct {
        Domain      string `json:"domain"`
        HasMX       bool   `json:"hasMX"`
        HasSPF      bool   `json:"hasSPF"`
        HasDMARC    bool   `json:"hasDMARC"`
        SPFRecord   string `json:"spfRecord"`
        DMARCRecord string `json:"dmarcRecord"`
}

// VerificationRequest represents the incoming request for domain verification
type VerificationRequest struct {
        Domain string `json:"domain"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
        Error string `json:"error"`
}

func main() {
        // Serve static files
        fs := http.FileServer(http.Dir("static"))
        http.Handle("/static/", http.StripPrefix("/static/", fs))

        // Define routes
        http.HandleFunc("/", handleIndex)
        http.HandleFunc("/verify", handleVerify)

        // Use port 8080 to avoid conflict with Flask
        port := "8080"
        
        // Check for environment variable override
        if envPort := os.Getenv("GO_PORT"); envPort != "" {
                port = envPort
        }

        // Start the server
        fmt.Printf("Go server running on http://0.0.0.0:%s\n", port)
        if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), nil); err != nil {
                log.Fatal("ListenAndServe: ", err)
        }
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
        // If path is not root, return 404
        if r.URL.Path != "/" {
                http.NotFound(w, r)
                return
        }

        // Serve index.html
        tmpl, err := template.ParseFiles(filepath.Join("templates", "index.html"))
        if err != nil {
                http.Error(w, "Error loading template", http.StatusInternalServerError)
                log.Println("Error loading template:", err)
                return
        }

        err = tmpl.Execute(w, nil)
        if err != nil {
                http.Error(w, "Error executing template", http.StatusInternalServerError)
                log.Println("Error executing template:", err)
                return
        }
}

func handleVerify(w http.ResponseWriter, r *http.Request) {
        // Only allow POST method
        if r.Method != http.MethodPost {
                w.Header().Set("Allow", http.MethodPost)
                respondWithError(w, "Method not allowed", http.StatusMethodNotAllowed)
                return
        }

        // Parse JSON request
        var req VerificationRequest
        decoder := json.NewDecoder(r.Body)
        if err := decoder.Decode(&req); err != nil {
                respondWithError(w, "Invalid request payload", http.StatusBadRequest)
                return
        }
        defer r.Body.Close()

        // Validate domain
        domain := strings.TrimSpace(req.Domain)
        if domain == "" {
                respondWithError(w, "Domain cannot be empty", http.StatusBadRequest)
                return
        }

        // Verify domain
        result, err := checkDomain(domain)
        if err != nil {
                respondWithError(w, fmt.Sprintf("Error verifying domain: %v", err), http.StatusInternalServerError)
                return
        }

        // Return result
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        if err := json.NewEncoder(w).Encode(result); err != nil {
                log.Printf("Error encoding response: %v", err)
        }
}

func respondWithError(w http.ResponseWriter, message string, statusCode int) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(statusCode)
        if err := json.NewEncoder(w).Encode(ErrorResponse{Error: message}); err != nil {
                log.Printf("Error encoding error response: %v", err)
        }
}

func checkDomain(domain string) (*VerificationResult, error) {
        var hasMX, hasSPF, hasDMARC bool
        var spfRecord, dmarcRecord string

        // Check for MX records
        mxRecords, err := net.LookupMX(domain)
        if err != nil {
                log.Printf("Error in MX lookup for %s: %v", domain, err)
                // Continue with other checks even if MX lookup fails
        } else if len(mxRecords) > 0 {
                hasMX = true
        }

        // Check for SPF records
        txtRecords, err := net.LookupTXT(domain)
        if err != nil {
                log.Printf("Error in TXT lookup for %s: %v", domain, err)
                // Continue with other checks even if TXT lookup fails
        } else {
                for _, record := range txtRecords {
                        if strings.HasPrefix(record, "v=spf1") {
                                hasSPF = true
                                spfRecord = record
                                break
                        }
                }
        }

        // Check for DMARC records
        dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
        if err != nil {
                log.Printf("Error in DMARC lookup for %s: %v", domain, err)
                // Continue with result even if DMARC lookup fails
        } else {
                for _, record := range dmarcRecords {
                        if strings.HasPrefix(record, "v=DMARC1") {
                                hasDMARC = true
                                dmarcRecord = record
                                break
                        }
                }
        }

        // Return verification result
        return &VerificationResult{
                Domain:      domain,
                HasMX:       hasMX,
                HasSPF:      hasSPF,
                HasDMARC:    hasDMARC,
                SPFRecord:   spfRecord,
                DMARCRecord: dmarcRecord,
        }, nil
}