const API_URL = 'https://api-nanolinq.ianchenn.com';

// Get DOM elements
const form = document.getElementById('statsForm');
const codeInput = document.getElementById('codeInput');
const statsResult = document.getElementById('statsResult');
const errorDiv = document.getElementById('error');
const errorMessage = document.getElementById('errorMessage');

// Stats elements
const originalUrl = document.getElementById('originalUrl');
const shortCode = document.getElementById('shortCode');
const clicks = document.getElementById('clicks');
const createdAt = document.getElementById('createdAt');
const lastAccessed = document.getElementById('lastAccessed');

// Handle form submit
form.addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const code = codeInput.value.trim();
    
    // Hide previous results/errors
    statsResult.classList.add('hidden');
    errorDiv.classList.add('hidden');
    
    try {
        const response = await fetch(`${API_URL}/${code}/stats`);
        const data = await response.json();
        
        if (response.ok) {
            // Success - show stats
            displayStats(data);
        } else {
            // Error from backend
            showError(data.error || 'Short code not found');
        }
    } catch (error) {
        // Network error
        showError('Network error. Please try again.');
        console.error('Error:', error);
    }
});

// Display stats
function displayStats(data) {
    originalUrl.href = data.url;
    originalUrl.textContent = data.url;
    shortCode.textContent = data.shortCode;
    clicks.textContent = data.clicks || 0;
    
    // Format dates
    createdAt.textContent = formatDate(data.createdAt);
    lastAccessed.textContent = data.lastAccessed 
        ? formatDate(data.lastAccessed) 
        : 'Never';
    
    statsResult.classList.remove('hidden');
}

// Format date to readable string
function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    });
}

// Show error message
function showError(message) {
    errorMessage.textContent = message;
    errorDiv.classList.remove('hidden');
}