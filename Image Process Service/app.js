const API_BASE = 'http://localhost:8080';
let authToken = localStorage.getItem('authToken');
let currentUser = localStorage.getItem('currentUser');
let currentImageId = null;

// Initialize
document.addEventListener('DOMContentLoaded', () => {
	if (authToken) {
		showDashboard();
	} else {
		showAuthPage();
	}

	setupEventListeners();
});

// Event Listeners
function setupEventListeners() {
	// Auth
	document.querySelectorAll('.auth-tab').forEach(btn => {
		btn.addEventListener('click', (e) => switchAuthTab(e.target.dataset.tab));
	});

	document.getElementById('loginForm').addEventListener('submit', (e) => {
		e.preventDefault();
		login();
	});

	document.getElementById('registerForm').addEventListener('submit', (e) => {
		e.preventDefault();
		register();
	});

	document.getElementById('logoutBtn').addEventListener('click', logout);

	// Navigation
	document.querySelectorAll('.nav-item').forEach(btn => {
		btn.addEventListener('click', (e) => switchSection(e.currentTarget.dataset.section));
	});

	// Upload area
	const uploadArea = document.getElementById('uploadArea');
	const fileInput = document.getElementById('fileInput');

	uploadArea.addEventListener('click', () => fileInput.click());

	uploadArea.addEventListener('dragover', (e) => {
		e.preventDefault();
		uploadArea.classList.add('dragover');
	});

	uploadArea.addEventListener('dragleave', () => {
		uploadArea.classList.remove('dragover');
	});

	uploadArea.addEventListener('drop', (e) => {
		e.preventDefault();
		uploadArea.classList.remove('dragover');
		if (e.dataTransfer.files.length) {
			fileInput.files = e.dataTransfer.files;
			uploadImage();
		}
	});

	fileInput.addEventListener('change', uploadImage);

	// Upload buttons
	const uploadBtnMain = document.getElementById('uploadBtnMain');
	if (uploadBtnMain) {
		uploadBtnMain.addEventListener('click', () => switchSection('upload'));
	}

	// Modal
	document.querySelectorAll('.modal-close').forEach(btn => {
		btn.addEventListener('click', (e) => {
			e.target.closest('.modal').classList.remove('active');
		});
	});

	document.querySelectorAll('.modal').forEach(modal => {
		modal.addEventListener('click', (e) => {
			if (e.target === modal) modal.classList.remove('active');
		});
	});

	// Transform
	document.getElementById('transformBtn').addEventListener('click', applyTransform);
	document.getElementById('deleteBtn').addEventListener('click', deleteImage);
}

// Auth Functions
function switchAuthTab(tab) {
	document.querySelectorAll('.auth-tab').forEach(btn => btn.classList.remove('active'));
	document.querySelectorAll('.auth-form').forEach(form => form.classList.remove('active'));

	document.querySelector(`[data-tab="${tab}"]`).classList.add('active');
	document.querySelector(`[data-form="${tab}"]`).classList.add('active');
}

// Navigation Functions
function switchSection(section) {
	document.querySelectorAll('.nav-item').forEach(btn => btn.classList.remove('active'));
	document.querySelectorAll('.section').forEach(s => s.classList.remove('active'));

	document.querySelector(`[data-section="${section}"]`).classList.add('active');
	document.getElementById(`${section}Section`).classList.add('active');
}

async function login() {
	const username = document.getElementById('loginUsername').value;
	const password = document.getElementById('loginPassword').value;
	const errorEl = document.getElementById('loginError');

	try {
		const res = await fetch(`${API_BASE}/login`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ username, password })
		});

		if (!res.ok) {
			errorEl.classList.add('show');
			errorEl.textContent = 'Invalid credentials';
			return;
		}

		const data = await res.json();
		authToken = data.jwt;
		currentUser = data.username;

		localStorage.setItem('authToken', authToken);
		localStorage.setItem('currentUser', currentUser);

		errorEl.classList.remove('show');
		showDashboard();
	} catch (err) {
		errorEl.classList.add('show');
		errorEl.textContent = 'Connection error';
	}
}

async function register() {
	const username = document.getElementById('regUsername').value;
	const password = document.getElementById('regPassword').value;
	const errorEl = document.getElementById('registerError');

	try {
		const res = await fetch(`${API_BASE}/register`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ username, password })
		});

		if (!res.ok) {
			errorEl.classList.add('show');
			errorEl.textContent = 'Username already exists';
			return;
		}

		const data = await res.json();
		authToken = data.jwt;
		currentUser = data.username;

		localStorage.setItem('authToken', authToken);
		localStorage.setItem('currentUser', currentUser);

		errorEl.classList.remove('show');
		showDashboard();
	} catch (err) {
		errorEl.classList.add('show');
		errorEl.textContent = 'Connection error';
	}
}

function logout() {
	authToken = null;
	currentUser = null;
	localStorage.removeItem('authToken');
	localStorage.removeItem('currentUser');
	showAuthPage();
}

// Page Navigation
function showAuthPage() {
	document.getElementById('authPage').classList.add('active');
	document.getElementById('dashboardPage').classList.remove('active');
}

function showDashboard() {
	document.getElementById('authPage').classList.remove('active');
	document.getElementById('dashboardPage').classList.add('active');
	document.getElementById('userDisplay').textContent = currentUser;
	loadImages();
}

// Images
async function loadImages() {
	try {
		const res = await fetch(`${API_BASE}/images?page=1&limit=100`, {
			headers: { 'Authorization': `Bearer ${authToken}` }
		});

		if (!res.ok) return;

		const data = await res.json();
		const gallery = document.getElementById('gallery');

		if (!data.images || data.images.length === 0) {
			gallery.innerHTML = `
				<div class="gallery-placeholder">
					<div class="placeholder-icon">📸</div>
					<p>No images yet</p>
					<p class="placeholder-hint">Upload your first image to get started</p>
				</div>
			`;
			return;
		}

		gallery.innerHTML = '';
		data.images.forEach(img => {
			const item = document.createElement('div');
			item.className = 'gallery-item';
			item.style.cursor = 'pointer';

			const image = document.createElement('img');
			image.className = 'gallery-item-image';
			image.src = img.url;
			image.alt = img.filename;

			const info = document.createElement('div');
			info.className = 'gallery-item-info';

			const name = document.createElement('div');
			name.className = 'gallery-item-name';
			name.textContent = img.filename;

			const size = document.createElement('div');
			size.className = 'gallery-item-size';
			size.textContent = `${img.width}×${img.height}`;

			info.appendChild(name);
			info.appendChild(size);
			item.appendChild(image);
			item.appendChild(info);

			item.addEventListener('click', () => openImageDetail(img.id));
			gallery.appendChild(item);
		});
	} catch (err) {
		console.error('Load error:', err);
	}
}

async function uploadImage() {
	const fileInput = document.getElementById('fileInput');
	const file = fileInput.files[0];
	const uploadError = document.getElementById('uploadError');
	const progressSection = document.getElementById('progressSection');

	if (!file) return;

	const formData = new FormData();
	formData.append('file', file);

	uploadError.classList.remove('show');
	progressSection.style.display = 'block';

	try {
		const xhr = new XMLHttpRequest();

		xhr.upload.addEventListener('progress', (e) => {
			const percent = (e.loaded / e.total) * 100;
			document.querySelector('.progress-fill').style.width = percent + '%';
			document.getElementById('progressText').textContent = Math.round(percent) + '%';
		});

		xhr.addEventListener('load', () => {
			if (xhr.status === 201) {
				progressSection.style.display = 'none';
				document.querySelector('.progress-fill').style.width = '0%';
				document.getElementById('progressText').textContent = '0%';
				fileInput.value = '';
				loadImages();
			} else {
				uploadError.classList.add('show');
				uploadError.textContent = 'Upload failed';
			}
		});

		xhr.addEventListener('error', () => {
			uploadError.classList.add('show');
			uploadError.textContent = 'Upload error';
		});

		xhr.open('POST', `${API_BASE}/images`);
		xhr.setRequestHeader('Authorization', `Bearer ${authToken}`);
		xhr.send(formData);
	} catch (err) {
		uploadError.classList.add('show');
		uploadError.textContent = 'Error';
	}
}

async function openImageDetail(imageId) {
	currentImageId = imageId;

	try {
		const res = await fetch(`${API_BASE}/images/${imageId}`, {
			headers: { 'Authorization': `Bearer ${authToken}` }
		});

		if (!res.ok) return;

		const img = await res.json();

		document.getElementById('detailImage').src = img.url;
		document.getElementById('modalImageName').textContent = img.filename;
		document.getElementById('imageDim').textContent = `${img.width}×${img.height}`;
		document.getElementById('imageFormat').textContent = img.format.toUpperCase();
		document.getElementById('imageSize').textContent = formatBytes(img.size);

		document.getElementById('resizeWidth').value = '';
		document.getElementById('resizeHeight').value = '';
		document.getElementById('rotateSelect').value = '0';
		document.getElementById('flipSelect').value = '';
		document.getElementById('formatSelect').value = 'jpeg';
		document.getElementById('grayscaleCheck').checked = false;
		document.getElementById('sepiaCheck').checked = false;

		document.getElementById('imageModal').classList.add('active');
	} catch (err) {
		console.error('Detail error:', err);
	}
}

async function applyTransform() {
	const width = parseInt(document.getElementById('resizeWidth').value) || 0;
	const height = parseInt(document.getElementById('resizeHeight').value) || 0;
	const rotate = parseInt(document.getElementById('rotateSelect').value);
	const flip = document.getElementById('flipSelect').value;
	const format = document.getElementById('formatSelect').value;
	const grayscale = document.getElementById('grayscaleCheck').checked;
	const sepia = document.getElementById('sepiaCheck').checked;
	const errorEl = document.getElementById('transformError');

	const transformations = {};

	if (width > 0 && height > 0) {
		transformations.resize = { width, height };
	}
	if (rotate !== 0) {
		transformations.rotate = rotate;
	}
	if (flip) {
		transformations.flip = flip;
	}
	if (grayscale || sepia) {
		transformations.filters = { grayscale, sepia };
	}
	transformations.format = format;

	try {
		const res = await fetch(`${API_BASE}/images/${currentImageId}/transform`, {
			method: 'POST',
			headers: {
				'Authorization': `Bearer ${authToken}`,
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ transformations })
		});

		if (!res.ok) {
			errorEl.classList.add('show');
			errorEl.textContent = 'Transform failed';
			return;
		}

		errorEl.classList.remove('show');
		alert('Image transformed successfully! Check your gallery for the new version.');
		document.getElementById('imageModal').classList.remove('active');
		loadImages();
	} catch (err) {
		errorEl.classList.add('show');
		errorEl.textContent = 'Error';
	}
}

async function deleteImage() {
	if (!confirm('Delete this image?')) return;

	try {
		const res = await fetch(`${API_BASE}/images/${currentImageId}`, {
			method: 'DELETE',
			headers: { 'Authorization': `Bearer ${authToken}` }
		});

		if (!res.ok) return;

		document.getElementById('imageModal').classList.remove('active');
		loadImages();
	} catch (err) {
		console.error('Delete error:', err);
	}
}

// Helpers
function formatBytes(bytes) {
	if (bytes === 0) return '0 Bytes';
	const k = 1024;
	const sizes = ['Bytes', 'KB', 'MB'];
	const i = Math.floor(Math.log(bytes) / Math.log(k));
	return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
}
