function closeModal() {
    const modal = document.getElementById('modal');
    const modalGramm = document.getElementById('modal__gramm');
    const searchInput = document.getElementById('searchInput');
    searchInput.value = '';
    modal.style.display = 'none';
    modalGramm.style.display = 'none';
}

function openModal(mealType) {
    const modal = document.getElementById('modal');
    const searchInput = document.getElementById('searchInput');
    searchInput.value = '';
    if (mealType !== undefined) {
        searchInput.setAttribute('data-mealType', mealType);
    }
    modal.style.display = 'block';
    document.querySelector('.close').addEventListener('click', closeModal);
    if (!mealType) {
        handleSearch({ target: searchInput });
    }
    console.log('Meal type in openModal:', mealType);
    console.log('Search input data-mealType attribute:', searchInput.getAttribute('data-mealType'));
}

function displaySearchResults(results, mealType) {
    console.log('Отображение результатов поиска:', results);
    const suggestionsList = document.getElementById('suggestions');
    console.log('Meal type in displaySearchResults:', mealType);
    if (!suggestionsList) {
        console.error("Элемент с id 'suggestions' не найден в DOM.");
        return;
    }

    suggestionsList.innerHTML = '';

    if (results.length === 0) {
        suggestionsList.style.display = 'none';
    } else {
        results.forEach(item => {
            const listItem = document.createElement('li');
            listItem.textContent = item.name;
            listItem.addEventListener('click', () => {
                document.getElementById('searchInput').value = item.name;
                suggestionsList.style.display = 'none';
                openGrammModal(item, mealType);
            });
            suggestionsList.appendChild(listItem);
        });
        suggestionsList.style.display = 'block';
    }
}

function updateSelectedProductInfo(selectedItem, mealType) {
    const selectedProductInfo = document.getElementById('selectedProductInfo');

    if (!selectedProductInfo) {
        console.error("Элемент с id 'selectedProductInfo' не найден в DOM.");
        return;
    }
    if (typeof mealType !== 'string' || mealType === undefined) {
        mealType = 'unknown';
    }

    selectedProductInfo.value = JSON.stringify({ ...selectedItem, mealType });
    console.log('Product data:', { ...selectedItem, mealType });
    console.log('Selected item:', selectedItem);
    console.log('Meal type updateSelectedProductInfo:', mealType);
}

function handleSearch(event) {
    console.log('handleSearch function called.');
    const searchInput = event.target;
    const mealType = searchInput.getAttribute('data-mealType');
    console.log('Meal type  handleSearch:', mealType);
    const searchInputValue = searchInput.value.trim().toLowerCase();
    console.log('Search input value:', searchInputValue);
    console.log('Search input data-mealType attribute:', searchInput.getAttribute('data-mealType'));

    const suggestionsList = document.getElementById('suggestions');
    if (searchInputValue === '') {
        suggestionsList.innerHTML = '';
        suggestionsList.style.display = 'none';
        console.log('Search input value is empty.');
        return;
    }

    console.log('Fetching data from server...');
    fetch('data.json')
        .then(response => response.json())
        .then(data => {
            console.log('Data received:', data);
            if (!suggestionsList) {
                console.error("Element with id 'suggestions' not found in DOM.");
                return;
            }
            const allItems = data.products.concat(data.dishes, data.drinks);
            const filteredData = allItems.filter(item => {
                return item.name.toLowerCase().includes(searchInputValue);
            });

            console.log('Filtered results:', filteredData);

            displaySearchResults(filteredData, mealType);

            console.log('displaySearchResults function called.');
            if (filteredData.length === 0) {
                console.log('Filtered results are empty.');
            } else {
                console.log('Filtered results are not empty.');
            }
        })
        .catch(error => {
            console.error('Error loading data.json file:', error);
            alert('Произошла ошибка при загрузке данных. Пожалуйста, попробуйте еще раз.');
        });
}

function openGrammModal(selectedItem, mealType) {
    if (!selectedItem) {
        console.error("Не передан объект selectedItem.");
        return;
    }
    const modalGramm = document.getElementById('modal__gramm');
    const cal = document.getElementById('caloriesValue');
    const pr = document.getElementById('proteinValue');
    const fat = document.getElementById('fatValue');
    const carb = document.getElementById('carbValue');
    modalGramm.style.display = 'block';
    document.getElementById('closeModalGramm').addEventListener('click', closeModal);
    caloriesInput.value = '';
    cal.textContent = '';
    pr.textContent = '';
    fat.textContent = '';
    carb.textContent = '';
    console.log('Product data:', { ...selectedItem, mealType });
    updateSelectedProductInfo(selectedItem, mealType);
    console.log('Meal type in openGrammModal:', mealType);
}
function updateTotalCaloriesDisplay(totalCalories) {
    const breakfastCaloriesDisplay = document.getElementById('count_of_calories_breakfast');
    const lunchCaloriesDisplay = document.getElementById('count_of_calories_lunch');
    const dinnerCaloriesDisplay = document.getElementById('count_of_calories_dinner');
    const totalCaloriesDisplay = document.getElementById('count_of_calories');
    const proteinDisplay = document.getElementById('count_of_protein');
    const fatDisplay = document.getElementById('count_of_fat');
    const carbDisplay = document.getElementById('count_of_carb');

    if (!totalCalories || !totalCalories.breakfast || !totalCalories.lunch || !totalCalories.dinner ||
        typeof totalCalories.breakfast.calories !== 'number' || typeof totalCalories.lunch.calories !== 'number' || typeof totalCalories.dinner.calories !== 'number') {
        console.error('Некорректные данные о калориях.');
        return;
    }

    const totalCaloriesSum = totalCalories.breakfast.calories + totalCalories.lunch.calories + totalCalories.dinner.calories;
    const totalProteinSum = totalCalories.breakfast.protein + totalCalories.lunch.protein + totalCalories.dinner.protein;
    const totalFatSum = totalCalories.breakfast.fat + totalCalories.lunch.fat + totalCalories.dinner.fat;
    const totalCarbohydrateSum = totalCalories.breakfast.carbohydrate + totalCalories.lunch.carbohydrate + totalCalories.dinner.carbohydrate;

    breakfastCaloriesDisplay.textContent = totalCalories.breakfast.calories.toFixed();
    lunchCaloriesDisplay.textContent = totalCalories.lunch.calories.toFixed();
    dinnerCaloriesDisplay.textContent = totalCalories.dinner.calories.toFixed();
    totalCaloriesDisplay.textContent = totalCaloriesSum.toFixed();
    proteinDisplay.textContent = totalProteinSum.toFixed(2);
    fatDisplay.textContent = totalFatSum.toFixed(2);
    carbDisplay.textContent = totalCarbohydrateSum.toFixed(2);

    localStorage.setItem('totalCalories', JSON.stringify(totalCalories));
}


function openProductsModal(mealType) {
    document.getElementById("modal__productBackdrop").style.display = "block";
    console.log('Meal type in openProductsModal:', mealType);
    fetch(`http://localhost:3000/products?mealtype=${mealType}&username=${localStorage.getItem('username')}`)
        .then(response => response.json())
        .then(data => {
            console.log('Data received from server:', data);

            const productList = document.getElementById("productList");
            productList.innerHTML = "";

            if (!data || !Array.isArray(data.Data)) {
                console.error('Некорректные данные с сервера:', data);
                return;
            }
            if (data.Data.length === 0) {
                document.getElementById("information").textContent = "Список продуктов пуст";
            } else {
                data.Data.forEach(product => {
                    const li = document.createElement("li");

                    const productInfo = document.createElement("div");
                    productInfo.classList.add("productInfo");

                    const productName = document.createElement("h3");
                    productName.textContent = `${product.name} - ${product.grams.toFixed(1)} г`;

                    const macronutrients = document.createElement("div");
                    macronutrients.textContent = `${product.calories.toFixed(0)} ккал, белки: ${product.protein.toFixed(1)} г, жиры: ${product.fat.toFixed(1)} г, угл.: ${product.carbohydrates.toFixed(1)} г`;
                    macronutrients.classList.add("macronutrients");

                    const deleteButton = document.createElement("button");
                    deleteButton.textContent = "Удалить";
                    deleteButton.classList.add("deleteButton");
                    deleteButton.onclick = function() {
                        deleteProduct(product.id, mealType);
                    };

                    productInfo.appendChild(productName);
                    productInfo.appendChild(macronutrients);
                    li.appendChild(productInfo);
                    li.appendChild(deleteButton);
                    productList.appendChild(li);
                });
                document.getElementById("information").textContent = "";
            }
        })
        .catch(error => console.error('Ошибка при загрузке продуктов:', error));
}






function deleteProduct(productId, mealType) {
    fetch(`http://localhost:3000/products?id=${productId}&username=${localStorage.getItem('username')}`, {
        method: 'DELETE',
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Ошибка удаления продукта');
        }
        return response.json();
    })
    .then(data => {
        console.log('Продукт успешно удален:', data);
        openProductsModal(mealType);
        updateTotalCaloriesDisplay(data.totalCalories);
    })
    .catch(error => console.error('Ошибка при удалении продукта:', error));
}


function closeProductsModal() {
    var modal = document.getElementById("modal__productBackdrop");
    const productList = document.getElementById("productList");
    productList.innerHTML = "";
    modal.style.display = "none";
}

function submitUsername() {
    const usernameInput = document.getElementById('usernameInput');
    const username = usernameInput.value.trim();
    if (username !== '') {
        localStorage.setItem('username', username);
        document.querySelector('.username__form').style.display = 'none';
        document.getElementById('main-content').style.display = 'block';
    } else {
        alert('Пожалуйста, введите имя пользователя.');
    }
}

window.addEventListener('load', function() {
    const savedUsername = localStorage.getItem('username');
    if (savedUsername) {
        const usernameInput = document.getElementById('usernameInput');
        usernameInput.value = savedUsername;
    }
});

document.addEventListener("DOMContentLoaded", function() {
    console.log('Событие DOMContentLoaded сработало.');
    const savedUsername = localStorage.getItem('username');
    function loadData() {
        fetch('data.json')
            .then(response => response.json())
            .then(data => {
                window.productsData = data.products;
            })
            .catch(error => {
                console.error('Error loading data.json file:', error);
            });
    }
    loadData();
    const savedTotalCalories = JSON.parse(localStorage.getItem('totalCalories'));
    
    if (savedTotalCalories) {
        updateTotalCaloriesDisplay(savedTotalCalories);
    } else {
        updateTotalCaloriesDisplay({ breakfast: { calories: 0, protein: 0, fat: 0, carbohydrate: 0 }, 
            lunch: { calories: 0, protein: 0, fat: 0, carbohydrate: 0 }, 
            dinner: { calories: 0, protein: 0, fat: 0, carbohydrate: 0 } });
    }

    const items = document.querySelectorAll('.services__item');
    const interval = 300;
    const observer = new IntersectionObserver(handleIntersection);

    items.forEach(item => {
        observer.observe(item);
    });

    console.log('IntersectionObserver initialized.');

    function handleIntersection(entries) {
        console.log('handleIntersection function called.');
        entries.forEach((entry, index) => {
            if (entry.isIntersecting) {
                setTimeout(() => {
                    entry.target.classList.add('show');
                    console.log('Element is intersecting and class "show" added.');
                }, index * interval);
            } else {
                entry.target.classList.remove('show');
                console.log('Element is not intersecting and class "show" removed.');
            }
        });
    }

    const searchInput = document.getElementById('searchInput');
    searchInput.addEventListener('input', handleSearch);

    searchInput.addEventListener('keydown', function(event) {
        if (event.key === 'Enter') {
            const selectedItem = JSON.parse(document.getElementById('selectedProductInfo').value);
            const mealType = searchInput.dataset.mealType;
            openGrammModal(selectedItem, mealType);
        }
    });

    const addButtons = document.querySelectorAll('.offer__btn');
    addButtons.forEach(button => {
        button.addEventListener('click', function(event) {
            event.preventDefault();
            const mealType = this.dataset.mealType;
            openModal(mealType);
        });
    });
    
    
    const openModalButton = document.querySelectorAll('.dest__btn');
    openModalButton.forEach(button => {
        const mealType = button.dataset.mealtype;
        button.addEventListener('click', function(event) {
            event.preventDefault();
            openProductsModal(mealType);
        });
    });
    function updateCalories() {
        const gramsInput = document.getElementById('caloriesInput');
        const grams = parseFloat(gramsInput.value);
        if (isNaN(grams)) return;
        const selectedItem = JSON.parse(document.getElementById('selectedProductInfo').value);
        if (!selectedItem) return;
        const caloriesPer100g = selectedItem.calories;
        const proteinPer100g = selectedItem.protein;
        const fatPer100g = selectedItem.fat;
        const carbPer100g = selectedItem.carbohydrates;

        const calories = (grams / 100) * caloriesPer100g;
        const protein = (grams / 100) * proteinPer100g;
        const fat = (grams / 100) * fatPer100g;
        const carbo = (grams / 100) * carbPer100g;

        const caloriesDisplay = document.getElementById('caloriesValue');
        const prDisplay = document.getElementById('proteinValue');
        const fatDisplay = document.getElementById('fatValue');
        const carboDisplay = document.getElementById('carbValue');

        caloriesDisplay.textContent = `${calories.toFixed()}`;
        prDisplay.textContent = `${protein === 0 ? '0' : protein.toFixed(1)}`;
        fatDisplay.textContent = `${fat === 0 ? '0' : fat.toFixed(1)}`;
        carboDisplay.textContent = `${carbo === 0 ? '0' : carbo.toFixed(1)}`;
    }


    const caloriesInput = document.getElementById('caloriesInput');

    caloriesInput.addEventListener('input', function() {
        const gramsValue = this.value.trim();

        if (gramsValue === '') {
            clearCaloriesInfo();
        } else {
            updateCalories();
        }
    });

    function clearCaloriesInfo() {
        const caloriesDisplay = document.getElementById('caloriesValue');
        const prDisplay = document.getElementById('proteinValue');
        const fatDisplay = document.getElementById('fatValue');
        const carboDisplay = document.getElementById('carbValue');
        caloriesDisplay.textContent = '';
        prDisplay.textContent = '';
        fatDisplay.textContent = '';
        carboDisplay.textContent = '';
    }

    function validateProductData(productData) {
        if (
            !productData ||
            typeof productData.name !== 'string' ||
            typeof productData.calories !== 'number' ||
            typeof productData.protein !== 'number' ||
            typeof productData.fat !== 'number' ||
            typeof productData.carbohydrates !== 'number' ||
            typeof productData.grams !== 'number' ||
            typeof productData.mealType !== 'string'
        ) {
            return false;
        }
        return true;
    }
    document.getElementById('confirmButton').addEventListener('click', function() {
        const gramsInput = document.getElementById('caloriesInput').value.trim();
        const grams = parseFloat(gramsInput);

        if (isNaN(grams)) {
            console.error('Некорректное количество грамм.');
            return;
        }

        const selectedProductInfo = JSON.parse(document.getElementById('selectedProductInfo').value);

        if (!selectedProductInfo) {
            console.error('Некорректные данные о продукте.');
            return;
        }

        const caloriesPer100g = selectedProductInfo.calories;
        const proteinPer100g = selectedProductInfo.protein;
        const fatPer100g = selectedProductInfo.fat;
        const carbPer100g = selectedProductInfo.carbohydrates;

        const calories = (grams / 100) * caloriesPer100g;
        const protein = (grams / 100) * proteinPer100g;
        const fat = (grams / 100) * fatPer100g;
        const carbo = (grams / 100) * carbPer100g;

        const productData = {
            name: selectedProductInfo.name,
            calories: calories,
            protein: protein,
            fat: fat,
            carbohydrates: carbo,
            grams: grams,
            mealType: selectedProductInfo.mealType
        };

        if (!validateProductData(productData)) {
            console.error('Некорректные данные о продукте.');
            return;
        }
        const url = `http://localhost:3000/products`;
        fetch(url, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({username:localStorage.getItem('username'), product:productData})
            })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Ошибка при сохранении данных о продукте на сервере.');
                }
                return response.json();
            })
            .then(data => {
                console.log('Данные о продукте успешно сохранены в базе данных:', data);
                updateTotalCaloriesDisplay(data.totalCalories);
                closeModal();
            })
            .catch(error => {
                console.error('Произошла ошибка при сохранении данных о продукте на сервере:', error);
            });
    });
    

    document.querySelector('.arrow-down').addEventListener('click', function(event) {
        event.preventDefault();
        const destItem = document.querySelector('.routes');
        destItem.scrollIntoView({ behavior: 'smooth' });
    });
});
