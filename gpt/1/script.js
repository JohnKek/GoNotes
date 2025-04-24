// Получаем элементы
const ball = document.querySelector('.ball');
const particlesContainer = document.getElementById('particles');

// Функция для создания частицы
function createParticle(x, y) {
    const particle = document.createElement('div');
    particle.classList.add('particle');
    particle.style.left = `${x}px`;
    particle.style.top = `${y}px`;
    particlesContainer.appendChild(particle);

    // Удаляем частицу после завершения анимации
    setTimeout(() => {
        particle.remove();
    }, 3000); // Время жизни частицы (1 секунда)
}

// Функция для обновления позиции шарика и создания частиц
function updateBallPosition(event) {
    const mouseX = event.clientX; // Координата X курсора
    const mouseY = event.clientY; // Координата Y курсора

    // Радиус шарика (половина его размера)
    const ballRadius = 25;

    // Устанавливаем новую позицию шарика
    ball.style.left = `${mouseX - ballRadius}px`;
    ball.style.top = `${mouseY - ballRadius}px`;

    // Создаем частицу в текущей позиции курсора
    createParticle(mouseX, mouseY);
}

// Слушаем событие движения мыши на документе
document.addEventListener('mousemove', updateBallPosition);