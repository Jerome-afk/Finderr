document.addEventListener('DOMContentLoaded', () => {
    const nowAiring = document.getElementById('new-releases');
    const popular = document.getElementById('popular');
    
    const moviePopular = 'https://api.themoviedb.org/3/movie/popular?api_key=22f697130e659aa6a1da22122fda7827&language=en-US';
    const movieAiring = 'https://api.themoviedb.org/3/movie/now_playing?api_key=22f697130e659aa6a1da22122fda7827&language=en-US';

    fetchMovies(moviePopular, movieAiring, nowAiring, popular);
});

function fetchMovies(apiUrl, apiUrl2, nowAiring, popular) {
    fetch(apiUrl2)
    .then(response => response.json())
    .then(data => {
        populateMediaRow(data.results, nowAiring);  // Now Airing
        addArrowFunctionality();
    })
    .catch(error => console.error('Error fetching movies:', error));

    fetch(apiUrl)
    .then(response => response.json())
    .then(data => {
        populateMediaRow(data.results, popular);  // Popular
        addArrowFunctionality();
    })
    .catch(error => console.error('Error fetching movies:', error));
}

function truncateString(text, wordLimit) {
    const words = text.split(' ');
    return words.length > wordLimit ? words.slice(0, wordLimit).join(' ') + '...' : text;
}

function populateMediaRow(mediaData, rowElement) {
    mediaData.forEach(movie => {
        const movieDiv = document.createElement('div');
        movieDiv.classList.add('movie-list-item');
        const truncateOverview = truncateString(movie.overview, 25);
        movieDiv.innerHTML = `
            <div class="media-description">
                <img src="https://image.tmdb.org/t/p/w500${movie.backdrop_path}" alt="${movie.title}" class="movie-list-item-img">
                <span class="movie-list-item-title">${movie.title}</span>
                <p class="movie-list-item-desc">${truncateOverview}</p>
                <button class="movie-list-item-button">Request</button>
            </div>
        `;
        rowElement.appendChild(movieDiv);
    });
}

function addArrowFunctionality() {
    const arrow = document.querySelectorAll('.arrow');
    const movieList = document.querySelectorAll('.movie-list');

    arrow.forEach((arrow, i) => {
        const itemNumber = movieList[i].querySelectorAll('.movie-list-item').length;
        let clickCounter = 0;
        let position = 0;

        arrow.addEventListener('click', () => {
            const ratio = Math.floor(window.innerWidth / 230);
            clickCounter++;
            position -= 230; // Move left by 230px on each click

            if (itemNumber - (5 + clickCounter) + (5 - ratio) < 0) {
                position = 0;
                clickCounter = 0;
                movieList[i].style.transform = "translateX(0)";
            } else {
                movieList[i].style.transform = `translateX(${position}px)`;
                movieList[i].style.transition = 'transform 0.5s ease';
            }
        });
    });
}

const ball = document.querySelector('.toggle-ball');
const items = document.querySelectorAll('.container, .movie-list-title, .navbar-container, .sidebar, .left-menu-icon, .toggle');

ball.addEventListener('click', () => {
    items.forEach(item => {
        item.classList.toggle('active');
    });
    ball.classList.toggle('active');
});