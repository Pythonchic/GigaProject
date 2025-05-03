document.addEventListener("DOMContentLoaded", function () {
    
    var priceValue = 0;
    // Обновляем отображение цены
    maxPrice.addEventListener("input", function () {
        priceValue.textContent = this.value;
    });

    // Фильтрация товаров
    applyFilterBtn.addEventListener("click", function () {
        const selectedSeason = seasonFilter.value.toLowerCase();
        let max = parseInt(maxPrice.value);
        let min = parseInt(minPrice.value);
        const minvalue = parseInt(lowerSlider.min);
        const maxvalue = parseInt(lowerSlider.max);


        if (min < minvalue) {
            min = minvalue;
            minPrice.value = minvalue;
        }

        if (max < minvalue) {
            max = maxvalue;
            maxPrice.value = maxvalue;
        }
        productCards.forEach(function (card) {
            const cost = parseInt(card.dataset.cost);
            const season = card.dataset.season.toLowerCase();

            const priceMatch = cost <= max && cost >= min;
            const seasonMatch = selectedSeason === "all" || season === selectedSeason;

            card.style.display = priceMatch && seasonMatch ? "block": "none";
        });
    });

    // Применяем фильтры сразу после загрузки
    applyFilterBtn.click();
});

// Форматирование числа с пробелами
function formatNumber(num) {
    return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, " ");
}

// Обновление позиции ползунков и трека
function updateSlider() {
    let lowerVal = parseInt(lowerSlider.value);
    let upperVal = parseInt(upperSlider.value);

    // Защита от пересечения значений
    if (lowerVal > upperVal) {
        [lowerVal, upperVal] = [upperVal, lowerVal];
        lowerSlider.value = lowerVal;
        upperSlider.value = upperVal;
    }

    // Расчет позиции трека
    const min = parseInt(lowerSlider.min);
    const max = parseInt(lowerSlider.max);
    const leftPercent = ((lowerVal - min) / (max - min)) * 100;
    const rightPercent = 100 - ((upperVal - min) / (max - min)) * 100;

    track.style.left = leftPercent + "%";
    track.style.right = rightPercent + "%";

    // Обновление полей ввода
    inputMin.value = lowerVal;
    inputMax.value = upperVal;
}

// Обработка изменения значений в полях ввода
function handleInputChange(e) {
    const input = e.target;
    let value = input.value === "" ? "" : parseInt(input.value);
    const minvalue = parseInt(lowerSlider.min);
    const maxvalue = parseInt(upperSlider.max);
    // Если поле пустое - не обновляем ползунок
    if (value === "") return;
    if (parseInt(value) < minvalue) return;

    if (parseInt(value) > maxvalue) {
        input.value = maxvalue;
    }
    const min = parseInt(input.min);
    const max = parseInt(input.max);

    // Корректировка значения
    value = Math.max(min, Math.min(max, value));

    // Обновление соответствующих ползунков
    if (input.id === "input-min") {
        lowerSlider.value = value;
        if (value > parseInt(upperSlider.value)) {
            upperSlider.value = value;
        }
    } else {
        upperSlider.value = value;
        if (value < parseInt(lowerSlider.value)) {
            lowerSlider.value = value;
        }
    }

    updateSlider();
}

// Обработка потери фокуса (если поле осталось пустым)
function handleInputBlur(e) {
    const input = e.target;
    if (input.value === "") {
        const defaultValue = input.id === "input-min" ? input.min : input.max;
        input.value = defaultValue;

        if (input.id === "input-min") {
            lowerSlider.value = defaultValue;
        } else {
            upperSlider.value = defaultValue;
        }

        updateSlider();
    }
}

const lowerSlider = document.getElementById("lower");
const upperSlider = document.getElementById("upper");
const inputMin = document.getElementById("input-min");
const inputMax = document.getElementById("input-max");
const track = document.getElementById("track");
const track2 = document.getElementById("track2");
const maxPrice = document.getElementById("input-max");
const minPrice = document.getElementById("input-min");
const seasonFilter = document.getElementById("season-filter");
const applyFilterBtn = document.getElementById("apply-filter");
const productCards = document.querySelectorAll(".product-card");

// Обработчики событий
lowerSlider.addEventListener("input", updateSlider);
upperSlider.addEventListener("input", updateSlider);

inputMin.addEventListener("input", handleInputChange);
inputMax.addEventListener("input", handleInputChange);

inputMin.addEventListener("blur", handleInputBlur);
inputMax.addEventListener("blur", handleInputBlur);

// Инициализация
updateSlider();
