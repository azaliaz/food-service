document.addEventListener('DOMContentLoaded', function() {
    const addWeightBtn = document.getElementById('addWeightBtn');
    const weightModal = document.getElementById('weightModal');
    const closeWeight = document.querySelector('.close-weight');
    const weightForm = document.getElementById('weightForm');
    const currentWeightDisplay = document.getElementById('currentWeightDisplay');
    const targetWeightDisplay = document.getElementById('targetWeightDisplay');
    const alertModal = document.getElementById('alertModal');
    const closeAlert = document.querySelector('.close-alert');
    const alertMessage = document.getElementById('alertMessage');
    const saveAnywayBtn = document.getElementById('saveAnywayBtn');

    let initialHeight, initialGender, initialWeight, currentWeight, targetWeight;
    
    function openWeightModal() {
        weightModal.style.display = 'block';
    }

    function closeWeightModal() {
        weightModal.style.display = 'none';
    }

    function openAlertModal(message) {
        alertMessage.textContent = message;
        alertModal.style.display = 'block';
    }

    function closeAlertModal() {
        alertModal.style.display = 'none';
    }

    addWeightBtn.addEventListener('click', function(event) {
        event.preventDefault(); 
        openWeightModal();
    });

    closeWeight.addEventListener('click', function() {
        closeWeightModal();
    });

   
    closeAlert.addEventListener('click', function() {
        closeAlertModal();
    });

    window.addEventListener('click', function(event) {
        if (event.target === weightModal) {
            closeWeightModal();
        } else if (event.target === alertModal) {
            closeAlertModal();
        }
    });

    weightForm.addEventListener('submit', function(event) {
        event.preventDefault(); 

        initialHeight = parseFloat(document.getElementById('initialHeight').value);
        initialGender = document.getElementById('initialGender').value;
        initialWeight = parseFloat(document.getElementById('initialWeight').value);
        currentWeight = parseFloat(document.getElementById('currentWeight').value);
        targetWeight = parseFloat(document.getElementById('targetWeight').value);

        let normalWeight;
        if (initialGender === 'male') {
            normalWeight = 0.9 * (initialHeight-100);
        } else if (initialGender === 'female') {
            normalWeight = 0.85 * (initialHeight-100);
        }

        if (targetWeight < normalWeight) {
            const message = `Ваш желанный вес (${targetWeight} кг) ниже нормы (${normalWeight.toFixed(2)} кг). Пожалуйста, пересмотрите свою цель.`;
            openAlertModal(message);
        } else {
            saveData();
            closeWeightModal();
        }

        saveAnywayBtn.addEventListener('click', function() {
            saveData();
            closeAlertModal();
            closeWeightModal();
        });
    });

    function saveData() {
        localStorage.setItem('initialHeight', initialHeight);
        localStorage.setItem('initialGender', initialGender);
        localStorage.setItem('initialWeight', initialWeight);
        localStorage.setItem('currentWeight', currentWeight);
        localStorage.setItem('targetWeight', targetWeight);

        
        currentWeightDisplay.textContent = ` ${currentWeight} кг`;
        targetWeightDisplay.textContent = ` ${targetWeight} кг`;
    }

   function initializeDisplay() {
        const storedCurrentWeight = localStorage.getItem('currentWeight');
        const storedTargetWeight = localStorage.getItem('targetWeight');

        if (storedCurrentWeight) {
            currentWeightDisplay.textContent = ` ${storedCurrentWeight} кг`;
        }

        if (storedTargetWeight) {
            targetWeightDisplay.textContent = ` ${storedTargetWeight} кг`;
        }
    }
    initializeDisplay();
});
