document.addEventListener('DOMContentLoaded', () => {
  const saveBtn = document.querySelector('.save-global-btn');
  const pdfWrapper = document.querySelector('.general-wrapper');
  let originalElements = [];

  const replaceInputsWithText = () => {
    originalElements = [];

    for (const el of document.querySelectorAll('input, textarea')) {
      const textEl = document.createElement('div');
      textEl.className = `text-replacement ${el.className}`;
      textEl.textContent = el.value || 'не указано';

      const style = window.getComputedStyle(el);
      Object.assign(textEl.style, {
        width: style.width,
        height: style.height,
        font: style.font,
        margin: style.margin,
        padding: style.padding,
        border: '1px solid transparent',
      });

      originalElements.push({
        original: el,
        parent: el.parentNode,
        nextSibling: el.nextSibling,
      });

      el.replaceWith(textEl);
    }
  };

  const restoreInputs = () => {
    for (const { original, parent, nextSibling } of originalElements) {
      if (nextSibling) {
        parent.insertBefore(original, nextSibling);
      } else {
        parent.appendChild(original);
      }
    }
    originalElements = [];
  };

  saveBtn.addEventListener('click', async () => {
    try {
      saveBtn.disabled = true;
      replaceInputsWithText();

      const canvas = await html2canvas(pdfWrapper, {
        scale: 2,
        useCORS: true,
        logging: false,
      });

      const pdf = new jspdf.jsPDF({
        orientation: 'p',
        unit: 'px',
        format: [canvas.width / 2, canvas.height / 2],
      });

      pdf.addImage(canvas, 'PNG', 0, 0, canvas.width / 2, canvas.height / 2);
      pdf.save('offer.pdf');
    } catch (err) {
      console.error('PDF Error:', err);
      alert(`Ошибка генерации: ${err.message}`);
    } finally {
      restoreInputs();
      saveBtn.disabled = false;
    }
  });
});

document.addEventListener('DOMContentLoaded', () => {
  const formElements = Array.from(document.querySelectorAll('input, textarea'));

  const handleKeyNavigation = (event) => {
    if (event.key === 'Enter' && event.target.tagName !== 'TEXTAREA') {
      event.preventDefault();
      const currentIndex = formElements.indexOf(event.target);

      if (currentIndex > -1 && currentIndex < formElements.length - 1) {
        formElements[currentIndex + 1].focus();
      }
    }
  };

  for (const el of formElements) {
    el.addEventListener('keydown', handleKeyNavigation);
  }
});

const updateFormElements = () => {
  formElements = Array.from(document.querySelectorAll('input, textarea'));
};

document.querySelector('.add-equipment').addEventListener('click', () => {
  addNewEquipmentField();
  updateFormElements();
});

if (event.key === 'Enter') {
  event.preventDefault();
  const currentIndex = formElements.indexOf(event.target);

  if (event.shiftKey) {
    if (currentIndex > 0) {
      formElements[currentIndex - 1].focus();
    }
  } else {
    if (currentIndex < formElements.length - 1) {
      formElements[currentIndex + 1].focus();
    }
  }
}

if (currentIndex > -1) {
  const nextElement = formElements[currentIndex + (event.shiftKey ? -1 : 1)];

  if (nextElement) {
    nextElement.focus();
    nextElement.scrollIntoView({
      behavior: 'smooth',
      block: 'center',
    });
  } else {
    if (!event.shiftKey) {
      document.querySelector('.save-global-btn').focus();
    }
  }
}
