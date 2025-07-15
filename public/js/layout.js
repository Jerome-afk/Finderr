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

    // Initialize carousel navigation
    const carousels = document.querySelectorAll('.media-scroller');
    
    carousels.forEach(carousel => {
        let isDown = false;
        let startX;
        let scrollLeft;

        carousel.addEventListener('mousedown', (e) => {
            isDown = true;
            startX = e.pageX - carousel.offsetLeft;
            scrollLeft = carousel.scrollLeft;
            carousel.style.cursor = 'grabbing';
        });

        carousel.addEventListener('mouseleave', () => {
            isDown = false;
            carousel.style.cursor = 'grab';
        });

        carousel.addEventListener('mouseup', () => {
            isDown = false;
            carousel.style.cursor = 'grab';
        });

        carousel.addEventListener('mousemove', (e) => {
            if(!isDown) return;
            e.preventDefault();
            const x = e.pageX - carousel.offsetLeft;
            const walk = (x - startX) * 2;
            carousel.scrollLeft = scrollLeft - walk;
        });
    });
});
