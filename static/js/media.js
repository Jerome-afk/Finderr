/*....................................................... Dark mode ...................................................... */
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
async function fetchPopularMovies() {
    const response = await fetch(`${BASE_URL}/movie/popular?api_key=${API_KEY}`);
    const data = await response.json();
    return data.results.slice(0, 10);
}

async function populateSlideShow() {
    const slideShow = document.getElementById('slideshow');
    const movies = await fetchPopularMovies();
    movies.forEach((movie, index) => {
        const slide = document.createElement('div');
        slide.classList.add('slide');
        if (index === 0) slide.classList.add('active');

        slide.innerHTML = `
        <img src="${IMAGE_BASE_URL + movie.backdrop_path}" alt="${movie.title}">
        <div class="slide-content">
            <h2>${movie.title}</h2>
            <p>${movie.overview}"</p>
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
async function fetchMoviesByCategory(category, containerId) {
    const response = await fetch(`${BASE_URL}/movie/${category}?api_key=${API_KEY}`);
    const data = await response.json();
    const container = document.getElementById(containerId);
  
    data.results.forEach((movie) => {
      const movieCard = document.createElement('div');
      movieCard.classList.add('movie-list-item');
      movieCard.innerHTML += `
        <div class="media-description">
            <div class="description>
                <img src="https://image.tmdb.org/t/p/w500${movie.backdrop_path}" alt="${movie.title}" class="movie-list-item-img">
                <div class="status-dot not-downloaded"></div>
                <button class="download-btn" onclick="downloadMovie('${movie.title}')><i class="fa-solid fa-download"></i></button>
            </div>
            <p class="title">${movie.title}<p>
        </div>
      `;
      container.appendChild(movieCard);
    });
  }

async function initializeContent() {
    await populateSlideShow();
    await fetchMoviesByCategory('popular', 'popular-category');
    await fetchMoviesByCategory('now_playing', 'now-playing-category');
    await fetchMoviesByCategory('top_rated', 'top-rated-category');
    await fetchMoviesByCategory('upcoming', 'upcoming-category');
}

initializeContent();