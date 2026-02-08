const API_URL = 'https://api-nanolinq.ianchenn.com';

// Get DOM elements
const form = document.getElementById('shortenForm');
const urlInput = document.getElementById('urlInput');
const resultDiv = document.getElementById('result');
const errorDiv = document.getElementById('error');
const shortUrlLink = document.getElementById('shortUrl');
const copyBtn = document.getElementById('copyBtn');
const errorMessage = document.getElementById('errorMessage');

// Handle form submit
form.addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const url = urlInput.value.trim();
    
    // Hide previous results/errors
    resultDiv.classList.add('hidden');
    errorDiv.classList.add('hidden');
    
    try {
        const response = await fetch(`${API_URL}/shorten`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ url }),
        });
        
        const data = await response.json();
        
        if (response.ok) {
            // Success - show result
            const shortUrl = `${API_URL}/${data.shortCode}`;
            shortUrlLink.href = shortUrl;
            shortUrlLink.textContent = shortUrl;
            resultDiv.classList.remove('hidden');
        } else {
            // Error from backend
            showError(data.error || 'Failed to create short URL');
        }
    } catch (error) {
        // Network error
        showError('Network error. Please try again.');
        console.error('Error:', error);
    }
});

// Handle copy button
copyBtn.addEventListener('click', () => {
    const url = shortUrlLink.href;
    navigator.clipboard.writeText(url).then(() => {
        // Visual feedback
        copyBtn.textContent = 'Copied!';
        setTimeout(() => {
            copyBtn.textContent = 'Copy';
        }, 2000);
    }).catch(err => {
        console.error('Failed to copy:', err);
    });
});

// Show error message
function showError(message) {
    errorMessage.textContent = message;
    errorDiv.classList.remove('hidden');
}