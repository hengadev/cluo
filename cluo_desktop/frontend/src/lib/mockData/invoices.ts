/**
 * Mock invoice data - Static, consistent across reloads.
 * Cases may have multiple invoices for different phases of work.
 */

import { INVOICE_LINE_ITEMS, PAYMENT_TERMS, CURRENCY, DATES, type DocumentStatus, type PaymentStatus } from './helpers';
import { CASE_IDS } from './cases';
import { CLIENT_IDS } from './clients';
import { CONTRACT_IDS } from './contracts';

// UUIDs for invoices
export const INVOICE_IDS = {
  // Draft case - no invoices yet

  // In progress cases
  inProgress1_inv1: 'd50e8400-e29b-41d4-a716-446655440101',

  // Ready cases - one invoice each (sent, not paid yet)
  ready1_inv1: 'd50e8400-e29b-41d4-a716-446655440201',
  ready2_inv1: 'd50e8400-e29b-41d4-a716-446655440202',
  ready3_inv1: 'd50e8400-e29b-41d4-a716-446655440203',

  // Released cases - multiple invoices each
  released1_inv1: 'd50e8400-e29b-41d4-a716-446655440301',
  released1_inv2: 'd50e8400-e29b-41d4-a716-446655440302',
  released2_inv1: 'd50e8400-e29b-41d4-a716-446655440303',
  released2_inv2: 'd50e8400-e29b-41d4-a716-446655440304',
  released3_inv1: 'd50e8400-e29b-41d4-a716-446655440305',
  released3_inv2: 'd50e8400-e29b-41d4-a716-446655440306',
  released4_inv1: 'd50e8400-e29b-41d4-a716-446655440307',
  released4_inv2: 'd50e8400-e29b-41d4-a716-446655440308',
  released4_inv3: 'd50e8400-e29b-41d4-a716-446655440309',
  released5_inv1: 'd50e8400-e29b-41d4-a716-446655440310',
  released5_inv2: 'd50e8400-e29b-41d4-a716-446655440311',
} as const;

export interface InvoiceLineItem {
  description: string;
  quantity: number;
  unitPrice: number;
  total: number;
}

export interface Invoice {
  id: string;
  caseId: string;
  clientId: string;
  invoiceNumber: string;
  issueDate: string;
  dueDate: string;
  lineItems: InvoiceLineItem[];
  totalAmount: number;
  taxRate: number;
  taxAmount: number;
  paymentStatus: PaymentStatus;
  paidAt: string | null;
  paidAmount: number | null;
  paymentMethod: string | null;
  linkedContractId: string | null;
  currency: string;
  paymentTerms: string;
  lateFee: number | null;
  lateFeeRate: number | null;
  status: DocumentStatus;
  createdAt: string;
  updatedAt: string;
}

export const invoices: Invoice[] = [
  // === IN PROGRESS CASES (draft invoices awaiting send) ===
  {
    id: INVOICE_IDS.inProgress1_inv1,
    caseId: CASE_IDS.inProgress1,
    clientId: CLIENT_IDS.person1,
    invoiceNumber: 'FAC-2025-001',
    issueDate: DATES.jan2024_30,
    dueDate: DATES.dec2025_01,
    lineItems: [
      { description: INVOICE_LINE_ITEMS[0].description, quantity: 1, unitPrice: 450, total: 450 },
      { description: INVOICE_LINE_ITEMS[3].description, quantity: 2, unitPrice: 150, total: 300 },
    ],
    totalAmount: 750,
    taxRate: 20,
    taxAmount: 150,
    paymentStatus: 'unpaid',
    paidAt: null,
    paidAmount: null,
    paymentMethod: null,
    linkedContractId: CONTRACT_IDS.inProgress1_contract1,
    currency: CURRENCY,
    paymentTerms: PAYMENT_TERMS[0],
    lateFee: null,
    lateFeeRate: null,
    status: 'draft',
    createdAt: DATES.jan2024_30,
    updatedAt: DATES.jan2024_30,
  },

  // === READY CASES (sent, awaiting payment) ===
  {
    id: INVOICE_IDS.ready1_inv1,
    caseId: CASE_IDS.ready1,
    clientId: CLIENT_IDS.person2,
    invoiceNumber: 'FAC-2024-028',
    issueDate: DATES.nov2024_05,
    dueDate: DATES.dec2025_01,
    lineItems: [
      { description: INVOICE_LINE_ITEMS[1].description, quantity: 3, unitPrice: INVOICE_LINE_ITEMS[1].unitPrice, total: 1140 },
      { description: INVOICE_LINE_ITEMS[6].description, quantity: 1, unitPrice: INVOICE_LINE_ITEMS[6].unitPrice, total: 120 },
    ],
    totalAmount: 1260,
    taxRate: 20,
    taxAmount: 252,
    paymentStatus: 'unpaid',
    paidAt: null,
    paidAmount: null,
    paymentMethod: null,
    linkedContractId: null,
    currency: CURRENCY,
    paymentTerms: PAYMENT_TERMS[0],
    lateFee: null,
    lateFeeRate: null,
    status: 'sent',
    createdAt: DATES.nov2024_05,
    updatedAt: DATES.nov2024_05,
  },
  {
    id: INVOICE_IDS.ready2_inv1,
    caseId: CASE_IDS.ready2,
    clientId: CLIENT_IDS.company2,
    invoiceNumber: 'FAC-2024-053',
    issueDate: DATES.nov2024_22,
    dueDate: DATES.jan2025_01,
    lineItems: [
      { description: INVOICE_LINE_ITEMS[0].description, quantity: 1, unitPrice: INVOICE_LINE_ITEMS[0].unitPrice, total: 450 },
      { description: INVOICE_LINE_ITEMS[4].description, quantity: 1, unitPrice: INVOICE_LINE_ITEMS[4].unitPrice, total: 220 },
      { description: INVOICE_LINE_ITEMS[5].description, quantity: 2, unitPrice: INVOICE_LINE_ITEMS[5].unitPrice, total: 360 },
    ],
    totalAmount: 1030,
    taxRate: 20,
    taxAmount: 206,
    paymentStatus: 'unpaid',
    paidAt: null,
    paidAmount: null,
    paymentMethod: null,
    linkedContractId: CONTRACT_IDS.ready2_contract1,
    currency: CURRENCY,
    paymentTerms: PAYMENT_TERMS[2],
    lateFee: null,
    lateFeeRate: null,
    status: 'sent',
    createdAt: DATES.nov2024_22,
    updatedAt: DATES.nov2024_22,
  },
  {
    id: INVOICE_IDS.ready3_inv1,
    caseId: CASE_IDS.ready3,
    clientId: CLIENT_IDS.lawyer2,
    invoiceNumber: 'FAC-2024-071',
    issueDate: DATES.dec2024_10,
    dueDate: DATES.feb2025_01,
    lineItems: [
      { description: INVOICE_LINE_ITEMS[5].description, quantity: 2, unitPrice: INVOICE_LINE_ITEMS[5].unitPrice, total: 360 },
      { description: INVOICE_LINE_ITEMS[11].description, quantity: 1, unitPrice: INVOICE_LINE_ITEMS[11].unitPrice, total: 250 },
      { description: INVOICE_LINE_ITEMS[4].description, quantity: 1, unitPrice: INVOICE_LINE_ITEMS[4].unitPrice, total: 220 },
    ],
    totalAmount: 830,
    taxRate: 20,
    taxAmount: 166,
    paymentStatus: 'unpaid',
    paidAt: null,
    paidAmount: null,
    paymentMethod: null,
    linkedContractId: null,
    currency: CURRENCY,
    paymentTerms: PAYMENT_TERMS[1],
    lateFee: null,
    lateFeeRate: null,
    status: 'sent',
    createdAt: DATES.dec2024_10,
    updatedAt: DATES.dec2024_10,
  },

  // === RELEASED CASES (all paid) ===
  {
    id: INVOICE_IDS.released1_inv1,
    caseId: CASE_IDS.released1,
    clientId: CLIENT_IDS.government1,
    invoiceNumber: 'FAC-2023-098',
    issueDate: DATES.jun2024_01,
    dueDate: DATES.jul2024_01,
    lineItems: [
      { description: INVOICE_LINE_ITEMS[0].description, quantity: 1, unitPrice: INVOICE_LINE_ITEMS[0].unitPrice, total: 450 },
      { description: INVOICE_LINE_ITEMS[3].description, quantity: 3, unitPrice: INVOICE_LINE_ITEMS[3].unitPrice, total: 450 },
    ],
    totalAmount: 900,
    taxRate: 20,
    taxAmount: 180,
    paymentStatus: 'paid',
    paidAt: DATES.jul2024_10,
    paidAmount: 1080,
    paymentMethod: 'Virement bancaire',
    linkedContractId: null,
    currency: CURRENCY,
    paymentTerms: PAYMENT_TERMS[0],
    lateFee: null,
    lateFeeRate: null,
    status: 'archived',
    createdAt: DATES.jun2024_01,
    updatedAt: DATES.aug2024_18,
  },
  {
    id: INVOICE_IDS.released1_inv2,
    caseId: CASE_IDS.released1,
    clientId: CLIENT_IDS.government1,
    invoiceNumber: 'FAC-2023-098-BIS',
    issueDate: DATES.jul2024_25,
    dueDate: DATES.aug2024_25,
    lineItems: [
      { description: INVOICE_LINE_ITEMS[4].description, quantity: 1, unitPrice: INVOICE_LINE_ITEMS[4].unitPrice, total: 220 },
      { description: INVOICE_LINE_ITEMS[5].description, quantity: 1, unitPrice: INVOICE_LINE_ITEMS[5].unitPrice, total: 180 },
    ],
    totalAmount: 400,
    taxRate: 20,
    taxAmount: 80,
    paymentStatus: 'paid',
    paidAt: DATES.aug2024_15,
    paidAmount: 480,
    paymentMethod: 'Virement bancaire',
    linkedContractId: null,
    currency: CURRENCY,
    paymentTerms: PAYMENT_TERMS[0],
    lateFee: null,
    lateFeeRate: null,
    status: 'archived',
    createdAt: DATES.jul2024_25,
    updatedAt: DATES.aug2024_18,
  },
  {
    id: INVOICE_IDS.released2_inv1,
    caseId: CASE_IDS.released2,
    clientId: CLIENT_IDS.person3,
    invoiceNumber: 'FAC-2023-121',
    issueDate: DATES.may2024_02,
    dueDate: DATES.jun2024_01,
    lineItems: [
      { description: INVOICE_LINE_ITEMS[1].description, quantity: 2, unitPrice: INVOICE_LINE_ITEMS[1].unitPrice, total: 760 },
      { description: INVOICE_LINE_ITEMS[6].description, quantity: 1, unitPrice: INVOICE_LINE_ITEMS[6].unitPrice, total: 120 },
    ],
    totalAmount: 880,
    taxRate: 20,
    taxAmount: 176,
    paymentStatus: 'paid',
    paidAt: DATES.jun2024_10,
    paidAmount: 1056,
    paymentMethod: 'Chèque',
    linkedContractId: null,
    currency: CURRENCY,
    paymentTerms: PAYMENT_TERMS[0],
    lateFee: null,
    lateFeeRate: null,
    status: 'archived',
    createdAt: DATES.may2024_02,
    updatedAt: DATES.sep2024_01,
  },
  {
    id: INVOICE_IDS.released2_inv2,
    caseId: CASE_IDS.released2,
    clientId: CLIENT_IDS.person3,
    invoiceNumber: 'FAC-2023-121-BIS',
    issueDate: DATES.jun2024_20,
    dueDate: DATES.jul2024_20,
    lineItems: [
      { description: INVOICE_LINE_ITEMS[4].description, quantity: 1, unitPrice: INVOICE_LINE_ITEMS[4].unitPrice, total: 220 },
      { description: INVOICE_LINE_ITEMS[7].description, quantity: 1, unitPrice: INVOICE_LINE_ITEMS[7].unitPrice, total: 120 },
    ],
    totalAmount: 340,
    taxRate: 20,
    taxAmount: 68,
    paymentStatus: 'paid',
    paidAt: DATES.jul2024_15,
    paidAmount: 408,
    paymentMethod: 'Chèque',
    linkedContractId: null,
    currency: CURRENCY,
    paymentTerms: PAYMENT_TERMS[0],
    lateFee: null,
    lateFeeRate: null,
    status: 'archived',
    createdAt: DATES.jun2024_20,
    updatedAt: DATES.sep2024_01,
  },
  {
    id: INVOICE_IDS.released3_inv1,
    caseId: CASE_IDS.released3,
    clientId: CLIENT_IDS.insurance1,
    invoiceNumber: 'FAC-2023-165',
    issueDate: DATES.may2024_15,
    dueDate: DATES.jun2024_15,
    lineItems: [
      { description: INVOICE_LINE_ITEMS[8].description, quantity: 1, unitPrice: INVOICE_LINE_ITEMS[8].unitPrice, total: 320 },
      { description: INVOICE_LINE_ITEMS[5].description, quantity: 1, unitPrice: INVOICE_LINE_ITEMS[5].unitPrice, total: 180 },
    ],
    totalAmount: 500,
    taxRate: 20,
    taxAmount: 100,
    paymentStatus: 'paid',
    paidAt: DATES.jun2024_10,
    paidAmount: 600,
    paymentMethod: 'Virement bancaire',
    linkedContractId: null,
    currency: CURRENCY,
    paymentTerms: PAYMENT_TERMS[1],
    lateFee: null,
    lateFeeRate: null,
    status: 'archived',
    createdAt: DATES.may2024_15,
    updatedAt: DATES.oct2024_08,
  },
  {
    id: INVOICE_IDS.released3_inv2,
    caseId: CASE_IDS.released3,
    clientId: CLIENT_IDS.insurance1,
    invoiceNumber: 'FAC-2023-165-BIS',
    issueDate: DATES.jul2024_10,
    dueDate: DATES.aug2024_25,
    lineItems: [
      { description: INVOICE_LINE_ITEMS[4].description, quantity: 1, unitPrice: INVOICE_LINE_ITEMS[4].unitPrice, total: 220 },
      { description: INVOICE_LINE_ITEMS[14].description, quantity: 1, unitPrice: INVOICE_LINE_ITEMS[14].unitPrice, total: 75 },
    ],
    totalAmount: 295,
    taxRate: 20,
    taxAmount: 59,
    paymentStatus: 'paid',
    paidAt: DATES.aug2024_20,
    paidAmount: 354,
    paymentMethod: 'Virement bancaire',
    linkedContractId: null,
    currency: CURRENCY,
    paymentTerms: PAYMENT_TERMS[1],
    lateFee: null,
    lateFeeRate: null,
    status: 'archived',
    createdAt: DATES.jul2024_10,
    updatedAt: DATES.oct2024_08,
  },
  {
    id: INVOICE_IDS.released4_inv1,
    caseId: CASE_IDS.released4,
    clientId: CLIENT_IDS.company1,
    invoiceNumber: 'FAC-2023-211',
    issueDate: DATES.may2024_02,
    dueDate: DATES.jun2024_01,
    lineItems: [
      { description: INVOICE_LINE_ITEMS[1].description, quantity: 2, unitPrice: INVOICE_LINE_ITEMS[1].unitPrice, total: 760 },
      { description: INVOICE_LINE_ITEMS[6].description, quantity: 1, unitPrice: INVOICE_LINE_ITEMS[6].unitPrice, total: 120 },
    ],
    totalAmount: 880,
    taxRate: 20,
    taxAmount: 176,
    paymentStatus: 'paid',
    paidAt: DATES.jun2024_05,
    paidAmount: 1056,
    paymentMethod: 'Virement bancaire',
    linkedContractId: CONTRACT_IDS.released4_contract1,
    currency: CURRENCY,
    paymentTerms: PAYMENT_TERMS[0],
    lateFee: null,
    lateFeeRate: null,
    status: 'archived',
    createdAt: DATES.may2024_02,
    updatedAt: DATES.nov2024_05,
  },
  {
    id: INVOICE_IDS.released4_inv2,
    caseId: CASE_IDS.released4,
    clientId: CLIENT_IDS.company1,
    invoiceNumber: 'FAC-2023-211-BIS',
    issueDate: DATES.jun2024_25,
    dueDate: DATES.jul2024_25,
    lineItems: [
      { description: INVOICE_LINE_ITEMS[1].description, quantity: 2, unitPrice: INVOICE_LINE_ITEMS[1].unitPrice, total: 760 },
      { description: INVOICE_LINE_ITEMS[7].description, quantity: 1, unitPrice: INVOICE_LINE_ITEMS[7].unitPrice, total: 120 },
    ],
    totalAmount: 880,
    taxRate: 20,
    taxAmount: 176,
    paymentStatus: 'paid',
    paidAt: DATES.jul2024_20,
    paidAmount: 1056,
    paymentMethod: 'Virement bancaire',
    linkedContractId: CONTRACT_IDS.released4_contract1,
    currency: CURRENCY,
    paymentTerms: PAYMENT_TERMS[0],
    lateFee: null,
    lateFeeRate: null,
    status: 'archived',
    createdAt: DATES.jun2024_25,
    updatedAt: DATES.nov2024_05,
  },
  {
    id: INVOICE_IDS.released4_inv3,
    caseId: CASE_IDS.released4,
    clientId: CLIENT_IDS.company1,
    invoiceNumber: 'FAC-2023-211-TER',
    issueDate: DATES.sep2024_01,
    dueDate: DATES.oct2024_01,
    lineItems: [
      { description: INVOICE_LINE_ITEMS[14].description, quantity: 1, unitPrice: INVOICE_LINE_ITEMS[14].unitPrice, total: 75 },
      { description: INVOICE_LINE_ITEMS[4].description, quantity: 1, unitPrice: INVOICE_LINE_ITEMS[4].unitPrice, total: 220 },
    ],
    totalAmount: 295,
    taxRate: 20,
    taxAmount: 59,
    paymentStatus: 'paid',
    paidAt: DATES.oct2024_05,
    paidAmount: 354,
    paymentMethod: 'Virement bancaire',
    linkedContractId: CONTRACT_IDS.released4_contract1,
    currency: CURRENCY,
    paymentTerms: PAYMENT_TERMS[0],
    lateFee: null,
    lateFeeRate: null,
    status: 'archived',
    createdAt: DATES.sep2024_01,
    updatedAt: DATES.nov2024_05,
  },
  {
    id: INVOICE_IDS.released5_inv1,
    caseId: CASE_IDS.released5,
    clientId: CLIENT_IDS.lawyer1,
    invoiceNumber: 'FAC-2023-255',
    issueDate: DATES.jun2024_01,
    dueDate: DATES.jul2024_01,
    lineItems: [
      { description: INVOICE_LINE_ITEMS[0].description, quantity: 1, unitPrice: INVOICE_LINE_ITEMS[0].unitPrice, total: 450 },
      { description: INVOICE_LINE_ITEMS[5].description, quantity: 1, unitPrice: INVOICE_LINE_ITEMS[5].unitPrice, total: 180 },
    ],
    totalAmount: 630,
    taxRate: 20,
    taxAmount: 126,
    paymentStatus: 'paid',
    paidAt: DATES.jul2024_10,
    paidAmount: 756,
    paymentMethod: 'Virement bancaire',
    linkedContractId: CONTRACT_IDS.released5_contract1,
    currency: CURRENCY,
    paymentTerms: PAYMENT_TERMS[1],
    lateFee: null,
    lateFeeRate: null,
    status: 'archived',
    createdAt: DATES.jun2024_01,
    updatedAt: DATES.dec2024_10,
  },
  {
    id: INVOICE_IDS.released5_inv2,
    caseId: CASE_IDS.released5,
    clientId: CLIENT_IDS.lawyer1,
    invoiceNumber: 'FAC-2023-255-BIS',
    issueDate: DATES.sep2024_12,
    dueDate: DATES.oct2024_25,
    lineItems: [
      { description: INVOICE_LINE_ITEMS[5].description, quantity: 1, unitPrice: INVOICE_LINE_ITEMS[5].unitPrice, total: 180 },
      { description: INVOICE_LINE_ITEMS[4].description, quantity: 1, unitPrice: INVOICE_LINE_ITEMS[4].unitPrice, total: 220 },
      { description: INVOICE_LINE_ITEMS[14].description, quantity: 1, unitPrice: INVOICE_LINE_ITEMS[14].unitPrice, total: 75 },
    ],
    totalAmount: 475,
    taxRate: 20,
    taxAmount: 95,
    paymentStatus: 'paid',
    paidAt: DATES.oct2024_20,
    paidAmount: 570,
    paymentMethod: 'Virement bancaire',
    linkedContractId: CONTRACT_IDS.released5_contract1,
    currency: CURRENCY,
    paymentTerms: PAYMENT_TERMS[1],
    lateFee: null,
    lateFeeRate: null,
    status: 'archived',
    createdAt: DATES.sep2024_12,
    updatedAt: DATES.dec2024_10,
  },
];

// Helper functions
export function getAllInvoices(): Invoice[] {
  return invoices;
}

export function getInvoiceById(id: string): Invoice | undefined {
  return invoices.find(i => i.id === id);
}

export function getInvoicesByCaseId(caseId: string): Invoice[] {
  return invoices.filter(i => i.caseId === caseId);
}

export function getInvoicesByClientId(clientId: string): Invoice[] {
  return invoices.filter(i => i.clientId === clientId);
}

export function getInvoicesByPaymentStatus(paymentStatus: PaymentStatus): Invoice[] {
  return invoices.filter(i => i.paymentStatus === paymentStatus);
}
