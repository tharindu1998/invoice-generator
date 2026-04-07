document.addEventListener('DOMContentLoaded', function () {
    // Set default dates
    const today = new Date().toISOString().split('T')[0];
    const due = new Date();
    due.setDate(due.getDate() + 30);

    const dateEl = document.getElementById('date');
    const dueDateEl = document.getElementById('due_date');
    if (dateEl) dateEl.value = today;
    if (dueDateEl) dueDateEl.value = due.toISOString().split('T')[0];

    // Amount calculation for initial row
    attachRowListeners(document.querySelector('.item-row'));

    // Phone validation
    initPhoneValidation('phone', 'phone_error');
    initPhoneValidation('customer_mobile', 'customer_mobile_error', true);

    // Form submit validation
    ['invoiceForm', 'customerForm'].forEach(function(id) {
        const form = document.getElementById(id);
        if (form) form.addEventListener('submit', function(e) {
            if (!validatePhone('phone', 'phone_error') | !validatePhone('customer_mobile', 'customer_mobile_error')) {
                e.preventDefault();
            }
        });
    });

    // Modal close
    const closeBtn = document.getElementById('closeModal');
    if (closeBtn) closeBtn.addEventListener('click', closeModal);

    window.addEventListener('click', function (e) {
        const modal = document.getElementById('previewModal');
        if (e.target === modal) closeModal();
    });
});

function initPhoneValidation(inputId, errorId, doLookup) {
    const inputEl = document.getElementById(inputId);
    if (!inputEl) return;

    let debounceTimer;
    inputEl.addEventListener('input', function () {
        // digits only
        this.value = this.value.replace(/\D/g, '');
        validatePhone(inputId, errorId);

        if (doLookup) {
            clearTimeout(debounceTimer);
            if (this.value.length >= 7) {
                debounceTimer = setTimeout(() => lookupCustomer(this.value), 400);
            }
        }
    });
}

function validatePhone(inputId, errorId) {
    const inputEl = document.getElementById(inputId);
    const errorEl = document.getElementById(errorId);
    if (!inputEl || !errorEl) return true;

    const val = inputEl.value.trim();
    if (!val) return true; // empty is handled by required attr

    if (val.length < 7 || val.length > 15) {
        errorEl.textContent = `Phone number must be between 7 and 15 digits (got ${val.length}).`;
        errorEl.classList.add('visible');
        inputEl.style.borderColor = 'var(--danger)';
        return false;
    }

    errorEl.classList.remove('visible');
    inputEl.style.borderColor = '';
    return true;
}

function lookupCustomer(phone) {
    fetch(`/api/customer?phone=${encodeURIComponent(phone)}`)
        .then(r => r.ok ? r.json() : null)
        .then(data => {
            if (!data) return;
            const nameEl    = document.getElementById('customer_name');
            const emailEl   = document.getElementById('customer_email');
            const addressEl = document.getElementById('customer_address');
            if (nameEl    && !nameEl.value)    nameEl.value    = data.name  || '';
            if (emailEl   && !emailEl.value)   emailEl.value   = data.email || '';
            if (addressEl && !addressEl.value) {
                const parts = [data.address_line1, data.address_line2].filter(Boolean);
                addressEl.value = parts.join(', ');
            }
        })
        .catch(() => {});
}

function attachRowListeners(row) {
    if (!row) return;
    row.querySelectorAll('input[name="item_quantity[]"], input[name="item_price[]"]')
        .forEach(input => input.addEventListener('input', () => recalcRow(row)));
}

function recalcRow(row) {
    const qty = parseFloat(row.querySelector('input[name="item_quantity[]"]').value) || 0;
    const price = parseFloat(row.querySelector('input[name="item_price[]"]').value) || 0;
    row.querySelector('input[name="item_amount[]"]').value = (qty * price).toFixed(2);
    updateTotal();
}

function updateTotal() {
    let total = 0;
    document.querySelectorAll('input[name="item_amount[]"]').forEach(i => {
        total += parseFloat(i.value) || 0;
    });
    const el = document.getElementById('grand-total');
    if (el) el.textContent = total.toFixed(2);
}

function addItem() {
    const tbody = document.getElementById('items-container');
    const row = document.createElement('tr');
    row.className = 'item-row';
    row.innerHTML = `
        <td><input type="text"   name="item_name[]"     required placeholder="Item description"></td>
        <td><input type="number" name="item_quantity[]" required min="1" value="1"></td>
        <td><input type="number" name="item_price[]"    required min="0" step="0.01" placeholder="0.00"></td>
        <td><input type="text"   name="item_amount[]"   readonly class="amount-display" placeholder="0.00"></td>
        <td>
            <button type="button" class="btn btn-danger btn-sm" onclick="removeItem(this)">
                <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                    <line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/>
                </svg>
            </button>
        </td>`;
    tbody.appendChild(row);
    attachRowListeners(row);
}

function removeItem(btn) {
    const rows = document.querySelectorAll('.item-row');
    if (rows.length <= 1) { alert('At least one item is required.'); return; }
    btn.closest('.item-row').remove();
    updateTotal();
}

function previewInvoice() {
    const form = document.getElementById('invoiceForm');
    if (!form) return;
    const fd = new FormData(form);

    const names = fd.getAll('item_name[]');
    const qtys  = fd.getAll('item_quantity[]');
    const prices = fd.getAll('item_price[]');

    let rows = '', total = 0;
    for (let i = 0; i < names.length; i++) {
        const amt = (parseFloat(qtys[i]) || 0) * (parseFloat(prices[i]) || 0);
        total += amt;
        rows += `<tr>
            <td>${names[i]}</td>
            <td style="text-align:center">${qtys[i]}</td>
            <td style="text-align:right">${parseFloat(prices[i]).toFixed(2)}</td>
            <td style="text-align:right">${amt.toFixed(2)}</td>
        </tr>`;
    }

    const bankName   = fd.get('bank_name')   || '';
    const bankAcc    = fd.get('bank_acc_no') || '';
    const bankBranch = fd.get('bank_branch') || '';
    const notes      = fd.get('notes')       || '';

    const paymentSection = (bankName || bankAcc || bankBranch || notes) ? `
        <div style="margin-top:1.25rem;padding:1rem;background:#f9fafb;border-radius:8px;font-size:0.875rem">
            <p style="font-size:0.75rem;font-weight:700;color:#6b7280;text-transform:uppercase;letter-spacing:0.5px;margin-bottom:0.75rem">Payment Details</p>
            <div style="display:grid;grid-template-columns:1fr 1fr;gap:0.5rem 1.5rem">
                ${bankName   ? `<p><span style="color:#6b7280">Bank</span><br><strong>${bankName}</strong></p>`         : ''}
                ${bankAcc    ? `<p><span style="color:#6b7280">Account No.</span><br><strong>${bankAcc}</strong></p>`   : ''}
                ${bankBranch ? `<p><span style="color:#6b7280">Branch</span><br><strong>${bankBranch}</strong></p>`     : ''}
            </div>
            ${notes ? `<p style="margin-top:0.75rem;color:#6b7280">Notes: ${notes}</p>` : ''}
        </div>` : '';

    document.getElementById('previewContent').innerHTML = `
        <p><strong>Customer:</strong> ${fd.get('customer_name') || '—'}</p>
        <p><strong>Mobile:</strong> ${fd.get('customer_mobile') || '—'}</p>
        <p style="margin-top:0.5rem"><strong>Date:</strong> ${fd.get('date')} &nbsp;|&nbsp; <strong>Due:</strong> ${fd.get('due_date')}</p>
        <table style="width:100%;border-collapse:collapse;margin:1.25rem 0;font-size:0.875rem">
            <thead>
                <tr style="background:#f9fafb;border-bottom:2px solid #e5e7eb">
                    <th style="padding:0.6rem 0.75rem;text-align:left">Description</th>
                    <th style="padding:0.6rem 0.75rem;text-align:center">Qty</th>
                    <th style="padding:0.6rem 0.75rem;text-align:right">Price</th>
                    <th style="padding:0.6rem 0.75rem;text-align:right">Amount</th>
                </tr>
            </thead>
            <tbody>${rows}</tbody>
            <tfoot>
                <tr style="border-top:2px solid #e5e7eb">
                    <td colspan="3" style="padding:0.75rem;text-align:right;font-weight:700">Total</td>
                    <td style="padding:0.75rem;text-align:right;font-weight:700">${total.toFixed(2)}</td>
                </tr>
            </tfoot>
        </table>
        ${paymentSection}`;

    document.getElementById('previewModal').classList.add('open');
}

function closeModal() {
    document.getElementById('previewModal').classList.remove('open');
}
