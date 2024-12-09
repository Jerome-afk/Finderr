// code for dark-mode
const htmlElement = document.documentElement;
const themeSwitch = document.querySelector('.switch input');

// Initialize theme based on localStorage
document.addEventListener('DOMContentLoaded', () => {
    const darkMode = localStorage.getItem('darkMode') === 'true';
    themeSwitch.checked = darkMode;
    if (darkMode) {
        htmlElement.classList.add('dark-mode');
    } else {
        htmlElement.classList.remove('dark-mode');
    }
});

//Toggle dark mode on checkbox change
themeSwitch.addEventListener('change', () => {
    if (themeSwitch.checked) {
        htmlElement.classList.add('dark-mode');
        localStorage.setItem('darkMode', 'true');
    } else {
        htmlElement.classList.remove('dark-mode');
        localStorage.setItem('darkMode', 'false');
    }
});

// category content
const API_KEY = '22f697130e659aa6a1da22122fda7827';
const BASE_URL = 'https://api.themoviedb.org/3';
const IMAGE_BASE_URL = 'https://image.tmdb.org/t/p/w780';

// Fetch popular movies
async function fetchPopularSeries() {
    const response = await fetch(`${BASE_URL}/tv/popular?api_key=${API_KEY}`);
    const data = await response.json();
    return data.results.slice(0, 10);
}

async function populateSlideShow() {
    const slideShow = document.getElementById('slideshow');
    const series = await fetchPopularSeries();
    series.forEach((tv, index) => {
        const slide = document.createElement('div');
        slide.classList.add('slide');
        if (index === 0) slide.classList.add('active');

        slide.innerHTML = `
        <img src="${IMAGE_BASE_URL + tv.backdrop_path}" alt="${tv.title}">
        <div class="slide-content">
            <h2>${tv.name}</h2>
            <p>${tv.overview}"</p>
            <button></button>
            <button>Info</button>
        </div>
        `;
        slideShow.appendChild(slide);
    });
    initializeSlideShow();
}

// Slideshow navigation
function initializeSlideShow() {
    const slides = document.querySelectorAll('.slide');
    let currentSlide = 0;

    const showSLide = (index) => {
        slides.forEach((slide, i) => {
            slide.classList.toggle('active', i === index);
        });
    };

    document.getElementById('prev-btn').addEventListener('click', () => {
        currentSlide = (currentSlide - 1 + slides.length) % slides.length;
        showSLide(currentSlide);
    });

    document.getElementById('next-btn').addEventListener('click', () => {
        currentSlide = (currentSlide + 1) % slides.length;
        showSLide(currentSlide);
    });
}

// Fetch movies for categories
async function fetchSeriesByCategory(category, containerId) {
    const response = await fetch(`${BASE_URL}/tv/${category}?api_key=${API_KEY}`);
    const data = await response.json();
    const container = document.getElementById(containerId);
  
    data.results.forEach((tv) => {
      const movieCard = document.createElement('div');
      movieCard.classList.add('movie-card');
      movieCard.innerHTML = `
        <img src="${IMAGE_BASE_URL + tv.poster_path}">
        <p>${tv.name}</p>
      `;
      container.appendChild(movieCard);
    });
  }

async function initializeContent() {
    await populateSlideShow();
    await fetchSeriesByCategory('popular', 'popular-category');
    await fetchSeriesByCategory('airing_today', 'airing-today-category');
    await fetchSeriesByCategory('top_rated', 'top-rated-category');
    await fetchSeriesByCategory('on_the_air', 'on-the-air-category');
}

initializeContent();