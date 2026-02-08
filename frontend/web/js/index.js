const API_URL = 'https://api-nanolinq.ianchenn.com';

// Tab elements
const tabShorten = document.getElementById('tabShorten');
const tabStats = document.getElementById('tabStats');
const contentShorten = document.getElementById('contentShorten');
const contentStats = document.getElementById('contentStats');

// Shorten form elements
const shortenForm = document.getElementById('shortenForm');
const urlInput = document.getElementById('urlInput');
const resultDiv = document.getElementById('result');
const errorDiv = document.getElementById('error');
const shortUrlLink = document.getElementById('shortUrl');
const displayCode = document.getElementById('displayCode');
const copyBtn = document.getElementById('copyBtn');
const viewStatsBtn = document.getElementById('viewStatsBtn');
const errorMessage = document.getElementById('errorMessage');

// Stats form elements
const statsForm = document.getElementById('statsForm');
const codeInput = document.getElementById('codeInput');
const statsResult = document.getElementById('statsResult');
const successStatsDiv = document.getElementById('successStats');
const errorStatsDiv = document.getElementById('errorStats');
const successStatsMessage = document.getElementById('successStatsMessage');
const errorStatsMessage = document.getElementById('errorStatsMessage');
const originalUrl = document.getElementById('originalUrl');
const shortCode = document.getElementById('shortCode');
const clicks = document.getElementById('clicks');
const createdAt = document.getElementById('createdAt');
const lastAccessed = document.getElementById('lastAccessed');
const deleteBtn = document.getElementById('deleteBtn');

// Store current short code
let currentShortCode = '';

// Store current short code being viewed
let viewingShortCode = '';

// Tab switching
tabShorten.addEventListener('click', () => {
    switchTab('shorten');
});

tabStats.addEventListener('click', () => {
    switchTab('stats');
});

function switchTab(tab) {
    if (tab === 'shorten') {
        tabShorten.classList.add('active');
        tabStats.classList.remove('active');
        contentShorten.classList.add('active');
        contentStats.classList.remove('active');
    } else {
        tabStats.classList.add('active');
        tabShorten.classList.remove('active');
        contentStats.classList.add('active');
        contentShorten.classList.remove('active');
    }
}

// Handle shorten form submit
shortenForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const url = urlInput.value.trim();
    
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
            // Store short code
            currentShortCode = data.shortCode;
            
            // Display results
            const shortUrl = `${API_URL}/${data.shortCode}`;
            shortUrlLink.href = shortUrl;
            shortUrlLink.textContent = shortUrl;
            displayCode.textContent = data.shortCode;
            
            resultDiv.classList.remove('hidden');
        } else {
            showError(data.error || 'Failed to create short URL');
        }
    } catch (error) {
        showError('Network error. Please try again.');
        console.error('Error:', error);
    }
});

// Handle copy button
copyBtn.addEventListener('click', () => {
    const url = shortUrlLink.href;
    navigator.clipboard.writeText(url).then(() => {
        copyBtn.textContent = 'Copied!';
        setTimeout(() => {
            copyBtn.textContent = 'Copy';
        }, 2000);
    }).catch(err => {
        console.error('Failed to copy:', err);
    });
});

// Handle view stats button
viewStatsBtn.addEventListener('click', () => {
    // Switch to stats tab
    switchTab('stats');
    
    // Auto-fill code input
    codeInput.value = currentShortCode;
    
    // Auto-submit (optional)
    statsForm.dispatchEvent(new Event('submit'));
});

// Handle stats form submit
statsForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const code = codeInput.value.trim();
    
    statsResult.classList.add('hidden');
    successStatsDiv.classList.add('hidden');
    errorStatsDiv.classList.add('hidden');
    
    try {
        const response = await fetch(`${API_URL}/${code}/stats`);
        const data = await response.json();
        
        if (response.ok) {
            viewingShortCode = code; // Store for delete
            displayStats(data);
        } else {
            showStatsError(data.error || 'Short code not found');
        }
    } catch (error) {
        showStatsError('Network error. Please try again.');
        console.error('Error:', error);
    }
});

// Handle delete button
deleteBtn.addEventListener('click', async () => {
    // Confirmation dialog
    const confirmed = confirm(
        `Are you sure you want to delete short code "${viewingShortCode}"?\n\n` +
        'This action cannot be undone.'
    );
    
    if (!confirmed) return;
    
    try {
        const response = await fetch(`${API_URL}/${viewingShortCode}`, {
            method: 'DELETE',
        });
        
        const data = await response.json();
        
        if (response.ok) {
            // Hide stats result
            statsResult.classList.add('hidden');
            errorStatsDiv.classList.add('hidden');
            
            // Show success message
            showStatsSuccess('Short code deleted successfully');
            
            // Clear input
            codeInput.value = '';
            viewingShortCode = '';
            
            // Optional: Switch back to Shorten tab after 2 seconds
            setTimeout(() => {
                switchTab('shorten');
                successStatsDiv.classList.add('hidden');
            }, 2000);
            
        } else {
            showStatsError(data.error || 'Failed to delete short code');
        }
    } catch (error) {
        showStatsError('Network error. Please try again.');
        console.error('Error:', error);
    }
});

// Display stats
function displayStats(data) {
    originalUrl.href = data.url;
    originalUrl.textContent = data.url;
    shortCode.textContent = data.shortCode;
    clicks.textContent = data.clicks || 0;
    createdAt.textContent = formatDate(data.createdAt);
    lastAccessed.textContent = data.lastAccessed 
        ? formatDate(data.lastAccessed) 
        : 'Never';
    
    statsResult.classList.remove('hidden');
}

// Format date
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

// Show errors
function showError(message) {
    errorMessage.textContent = message;
    errorDiv.classList.remove('hidden');
}

// Show success message
function showStatsSuccess(message) {
    successStatsMessage.textContent = message;
    successStatsDiv.classList.remove('hidden');
}

function showStatsError(message) {
    errorStatsMessage.textContent = message;
    errorStatsDiv.classList.remove('hidden');
}