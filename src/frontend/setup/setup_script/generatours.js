document.addEventListener('DOMContentLoaded', () => {
  const formElements = document.querySelectorAll('input, textarea');
  formElements.forEach((el, index) => {
    el.addEventListener('keydown', (event) => {
      if (event.key === 'Enter' && el.tagName !== 'TEXTAREA') {
        event.preventDefault();
        const nextEl = formElements[index + 1];
        if (nextEl) {
          nextEl.focus();
        }
      }
    });
  });

  const saveBtn = document.querySelector('.save-global-btn');
  const pdfWrapper = document.querySelector('.general-wrapper');

  saveBtn.addEventListener('click', async () => {
    try {

      const canvas = await html2canvas(pdfWrapper, { scale: 2 });
      const imgData = canvas.toDataURL('image/png');

      const pdf = new jsPDF('p', 'pt', 'a4');
      const pageWidth = pdf.internal.pageSize.getWidth();

      const imgWidth = pageWidth;
      const imgHeight = (canvas.height * imgWidth) / canvas.width;

      pdf.addImage(imgData, 'PNG', 0, 0, imgWidth, imgHeight);
      pdf.save('document.pdf');
    } catch (err) {
      console.error('Ошибка генерации PDF:', err);
      alert('Не удалось сгенерировать PDF');
    }
  });
});
