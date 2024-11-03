class StarryBackground {
    constructor(container, options = {}) {
        this.container = typeof container === 'string' ? 
            document.querySelector(container) : container;
        
        this.options = {
            starCount: options.starCount || 50,
            minDuration: options.minDuration || 2,
            maxDuration: options.maxDuration || 5
        };

        this.stars = [];
        this.init();
    }

    init() {
        // Create stars
        for (let i = 0; i < this.options.starCount; i++) {
            this.createStar();
        }

        // Handle resize
        window.addEventListener('resize', () => this.handleResize());
    }

    createStar() {
        const star = document.createElement('div');
        star.className = 'star';
        
        // Random position
        star.style.left = `${Math.random() * 100}%`;
        star.style.top = `${Math.random() * 100}%`;
        
        // Random animation duration
        const duration = this.options.minDuration + 
            Math.random() * (this.options.maxDuration - this.options.minDuration);
        star.style.setProperty('--duration', `${duration}s`);

        this.container.appendChild(star);
        this.stars.push(star);
    }

    handleResize() {
        // Optional: Reposition stars on resize
        this.stars.forEach(star => {
            star.style.left = `${Math.random() * 100}%`;
            star.style.top = `${Math.random() * 100}%`;
        });
    }

    destroy() {
        this.stars.forEach(star => star.remove());
        this.stars = [];
        window.removeEventListener('resize', this.handleResize);
    }
}