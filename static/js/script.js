document.addEventListener('DOMContentLoaded', function() {
    const verificationForm = document.getElementById('verificationForm');
    const domainInput = document.getElementById('domain');
    const verifyButton = document.getElementById('verifyButton');
    const spinnerIcon = document.getElementById('spinnerIcon');
    const loadingMessage = document.getElementById('loadingMessage');
    const errorMessage = document.getElementById('errorMessage');
    const errorContent = document.getElementById('errorContent');
    const resultContainer = document.getElementById('resultContainer');
    const resultDomain = document.getElementById('resultDomain');
    const mxStatus = document.getElementById('mxStatus');
    const spfStatus = document.getElementById('spfStatus');
    const dmarcStatus = document.getElementById('dmarcStatus');
    const spfRecord = document.getElementById('spfRecord');
    const dmarcRecord = document.getElementById('dmarcRecord');

    // Function to show loading state
    function showLoading() {
        spinnerIcon.classList.remove('d-none');
        loadingMessage.classList.remove('d-none');
        errorMessage.classList.add('d-none');
        resultContainer.classList.add('d-none');
        verifyButton.disabled = true;
    }

    // Function to hide loading state
    function hideLoading() {
        spinnerIcon.classList.add('d-none');
        loadingMessage.classList.add('d-none');
        verifyButton.disabled = false;
    }

    // Function to show error
    function showError(message) {
        errorMessage.classList.remove('d-none');
        errorContent.textContent = message;
        resultContainer.classList.add('d-none');
    }

    // Function to create status badge
    function createStatusBadge(status, text) {
        let badgeClass = status ? 'bg-success' : 'bg-danger';
        let badgeText = status ? 'Configured' : 'Not Configured';
        
        if (text) {
            badgeText = text;
        }
        
        return `<span class="badge ${badgeClass}">${badgeText}</span>`;
    }

    // Function to display verification results
    function displayResults(data) {
        resultContainer.classList.remove('d-none');
        resultDomain.textContent = data.domain;
        
        // Display MX status
        mxStatus.innerHTML = createStatusBadge(data.hasMX);
        
        // Display SPF status
        spfStatus.innerHTML = createStatusBadge(data.hasSPF);
        
        // Display DMARC status
        dmarcStatus.innerHTML = createStatusBadge(data.hasDMARC);
        
        // Display record details
        if (data.spfRecord) {
            spfRecord.innerHTML = `<code>${data.spfRecord}</code>`;
        } else {
            spfRecord.innerHTML = `<code>No SPF record found</code>`;
        }
        
        if (data.dmarcRecord) {
            dmarcRecord.innerHTML = `<code>${data.dmarcRecord}</code>`;
        } else {
            dmarcRecord.innerHTML = `<code>No DMARC record found</code>`;
        }
    }

    // Form submission handler
    verificationForm.addEventListener('submit', function(e) {
        e.preventDefault();
        
        const domain = domainInput.value.trim();
        
        if (!domain) {
            showError('Please enter a domain name');
            return;
        }
        
        // Start the verification process
        showLoading();
        
        fetch('/verify', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ domain: domain }),
        })
        .then(response => {
            if (!response.ok) {
                return response.json().then(data => {
                    throw new Error(data.error || 'Verification failed');
                });
            }
            return response.json();
        })
        .then(data => {
            hideLoading();
            displayResults(data);
        })
        .catch(error => {
            hideLoading();
            showError(error.message || 'An unexpected error occurred');
        });
    });
});