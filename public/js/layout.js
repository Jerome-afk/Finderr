document.addEventListener("DOMContentLoaded", function() {
    // Toggle the sidebar
    const sidebar = document.querySelector(".sidebar");
    if (sidebar) {
        // Keep it minimised by default
        sidebar.classList.add("minimised");

        // Toggle on click
        sidebar.addEventListener("click", function() {
            this.classList.toggle("minimised");
        });
    }

    // Carousel functionality
    const carousel = {
        slides: document.querySelectorAll(".carousel-slide"),
        indicators: document.querySelectorAll(".indicator"),
        prevBtn: document.querySelector(".control-prev"),
        nextBtn: document.querySelector(".control-next"),
        currentIndex: 0,
        interval: null,
        intervalTime: 5000, // 5 seconds

        init() {
            // Set up event listeners
            this.prevBtn.addEventListener("click", () => this.prevSlide());
            this.nextBtn.addEventListener("click", () => this.nextSlide());

            // Indicator click
            this.indicators.forEach(indicator => {
                indicator.addEventListener('click', () => {
                    this.goToSlide(parseInt(indicator.dataset.index));
                });
            });

            // Start the carousel
            this.startAutoRotation();

            // Pause on hover
            const carouselContainer = document.querySelector(".carousel-container");
            carouselContainer.addEventListener("mouseenter", () => {
                this.pauseAutoRotation();
            });
            carouselContainer.addEventListener("mouseleave", () => {
                this.startAutoRotation();
            });
        },

        updateSlide() {
            // Hide all slides
            this.slides.forEach(slide => slide.classList.remove("active"));
            this.indicators.forEach(indicator => indicator.classList.remove("active"));

            // Show the current slide
            this.slides[this.currentIndex].classList.add("active");
            this.indicators[this.currentIndex].classList.add("active");
        },

        nextSlide() {
            this.currentIndex = (this.currentIndex + 1) % this.slides.length;
            this.updateSlide();
            this.resetAutoRotation();
        },

        prevSlide() {
            this.currentIndex = (this.currentIndex - 1 + this.slides.length) % this.slides.length;
            this.updateSlide();
            this.resetAutoRotation();
        },

        goToSlide(index) {
            this.currentIndex = index;
            this.updateSlide();
            this.resetAutoRotation();
        },

        startAutoRotation() {
            this.interval = setInterval(() => this.nextSlide(), this.intervalTime);
        },

        pauseAutoRotation() {
            clearInterval(this.interval);
        },

        resetAutoRotation() {
            this.pauseAutoRotation();
            this.startAutoRotation();
        }
    };

    // Initialize the carousel
    carousel.init();

    // Format the backdrop to load with highest quality
    const lazyLoadBackdrops = () => {
        const slides = document.querySelectorAll('.carousel-slude:not(.active)');
        slides.forEach(slide => {
            const bg = slide.querySelector('.slide-background');
            if (bg) {
                const currentSrc = bg.style.backgroundImage.match(/url\(["']?(.*?)["']?\)/);
                const highResSrc = currentSrc.replace('w500', 'original');
                bg.style.backgroundImage = `url(${highResSrc})`;
            }
        });
    };

    // Load higher quality images after initial load
    setTimeout(() => {
        lazyLoadBackdrops();
    }, 1000);
});
