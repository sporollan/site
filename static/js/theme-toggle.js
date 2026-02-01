// Theme toggle functionality
(function() {
    const storageKey = 'theme';
    const lightClass = 'theme-light';
    const darkClass = 'theme-dark';
    
    // Get current theme from localStorage or system preference
    function getThemePreference() {
        const saved = localStorage.getItem(storageKey);
        if (saved === 'light' || saved === 'dark') {
            return saved;
        }
        // Use system preference
        return window.matchMedia('(prefers-color-scheme: light)').matches ? 'light' : 'dark';
    }
    
    // Set theme on html element
    function setTheme(theme) {
        const html = document.documentElement;
        html.classList.remove(lightClass, darkClass);
        html.classList.add(theme === 'light' ? lightClass : darkClass);
        localStorage.setItem(storageKey, theme);
        updateIcon(theme);
    }
    
    // Toggle between light and dark
    function toggleTheme() {
        const current = getThemePreference();
        const newTheme = current === 'light' ? 'dark' : 'light';
        setTheme(newTheme);
    }
    
    // Update icon visibility
    function updateIcon(theme) {
        const lightIcon = document.getElementById('theme-icon-light');
        const darkIcon = document.getElementById('theme-icon-dark');
        if (!lightIcon || !darkIcon) return;
        
        if (theme === 'light') {
            lightIcon.style.display = 'none';
            darkIcon.style.display = 'block';
        } else {
            lightIcon.style.display = 'block';
            darkIcon.style.display = 'none';
        }
    }
    
    // Initialize theme on page load
    function initTheme() {
        const theme = getThemePreference();
        setTheme(theme);
        
        // Add click event to toggle button
        const toggleBtn = document.getElementById('theme-toggle');
        if (toggleBtn) {
            toggleBtn.addEventListener('click', toggleTheme);
        }
    }
    
    // Run when DOM is ready
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', initTheme);
    } else {
        initTheme();
    }
})();