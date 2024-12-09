const API_KEY = '2';
const BASE_URL = 'https://api.deezer.com';
const sidebar = document.querySelector('.sidebar');
const menuButton = document.querySelector('.menu-btn');
const closeButton = document.querySelector('.close-btn');
const username = "";
const profileIcon = document.querySelector('.profile-icon');
const notifSidebar = document.querySelector('.notification-sidebar');
const notifIcon = document.querySelector('.notif-icon');
const closeNotifButton = document.querySelector('.close-notif-btn');

if (username) {
    profileIcon.src = `/users/${username}/profile.jpeg`;
} else {
    profileIcon.src = "/static/images/noprofile.png";
}

// Toggle sidebar visibility
menuButton.addEventListener('click', () => {
    sidebar.classList.add('open');
});

// Close sidebar when close button is clicked
closeButton.addEventListener('click', () => {
    sidebar.classList.remove('open');
});

// Close sidebar when clicking outside it
document.addEventListener('click', (event) => {
    if (!sidebar.contains(event.target) && !menuButton.contains(event.target)) {
        sidebar.classList.remove('open');
    }
});

// Toggle Notification Sidebar
notifIcon.addEventListener('click', () => {
    notifSidebar.classList.add('open');
});

// Close Notification Sidebar
closeNotifButton.addEventListener('click', () => {
    notifSidebar.classList.remove('open');
});

// Close Notification Sidebar on Outside Click
document.addEventListener('click', (event) => {
    if (
        !notifSidebar.contains(event.target) &&
        !notifIcon.contains(event.target)
    ) {
        notifSidebar.classList.remove('open');
    }
});

// Fetch and display trending artists
async function fetchTrendingArtists() {
    try {
        const response = await fetch(`${BASE_URL}/chart/0/artists`);
        const data = await response.json();
        const artists = data.data || [];
        const container = document.querySelector('#trending-artists .artists');
        container.innerHTML = ''; // Clear container before appending new content

        artists.forEach(artist => {
            const div = document.createElement('div');
            div.innerHTML = `
                <img src="${artist.picture_medium || ''}" alt="${artist.name}" />
                <p>${artist.name}</p>
            `;
            container.appendChild(div);
        });
    } catch (error) {
        console.error('Error fetching trending artists:', error);
    }
}

// Fetch popular albums
async function fetchPopularAlbums() {
    try {
        const response = await fetch(`${BASE_URL}/chart/0/albums`);
        const data = await response.json();
        const albums = data.data || [];
        const container = document.querySelector('#popular-albums .albums');
        container.innerHTML = ''; // Clear container before appending new content

        albums.forEach(album => {
            const div = document.createElement('div');
            div.innerHTML = `
                <img src="${album.cover_medium || ''}" alt="${album.title}" />
                <p>${album.title} by ${album.artist.name}</p>
            `;
            container.appendChild(div);
        });
    } catch (error) {
        console.error('Error fetching popular albums:', error);
    }
}

// Fetch new releases
async function fetchNewReleases() {
    try {
        const response = await fetch(`${BASE_URL}/chart/0/tracks`);
        const data = await response.json();
        const releases = data.data || [];
        const container = document.querySelector('#new-releases .releases');
        container.innerHTML = ''; // Clear container before appending new content

        releases.forEach(release => {
            const div = document.createElement('div');
            div.innerHTML = `
                <img src="${release.album.cover_medium || ''}" alt="${release.title}" />
                <p>${release.title} by ${release.artist.name}</p>
            `;
            container.appendChild(div);
        });
    } catch (error) {
        console.error('Error fetching new releases:', error);
    }
}

// Initialize the fetching process
async function initializeMusicData() {
    await fetchTrendingArtists();
    await fetchPopularAlbums();
    await fetchNewReleases();
}

// Call the initialization function
initializeMusicData();