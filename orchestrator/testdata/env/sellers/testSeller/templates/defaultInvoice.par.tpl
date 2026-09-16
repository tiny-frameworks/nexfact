
# Dateiname: ./test/data/GSF_DIN-5008_Invoice-Max.par
# Demo und Test Paraemter File. 
# Zugehöriges .ott File: ./templates/GSF_DIN-5008_Invoice-Max.ott

# Dateien - Ablage (Globals)
# Präfix: g-
#------------------
# ott: Vorlage zum Erzeugen des Dokuments
# odt: Speicherort des fertigen ODF-Dokuments (saveTo)
# pdf: Speicherort des fertigen PDF-Dokuments (exportTo)
{{ with .Invoice.Paths }}
g-ott = {{ .OTTfile }}
g-odt = {{ .ODTfile }}
g-pdf = {{ .PDFfile }}
{{ end }}

# einfache Parameter
# Präfix: p-
#-------------------
# Absender - Adresse


# Postvermerkzeilen / Absender - AdresseZeile
{{ with .Invoice.Headers }}
p-senderLine = {{ .SenderLine }}
p-postalNote1 = {{ .PostalNote1 }}
p-postalNote2 = {{ .PostalNote2 }}
p-postalNote3 = {{ .PostalNote3 }}
p-infoDepartment = {{ .Department }}
p-infoDate = {{ .InfoDate | toPARdate }}
{{ end }}

# Empfaengerzeilen
{{ with .Invoice.Invoice.Buyer }}
p-recipient1 = {{ .Name }}
p-recipient2 = {{ .Name2 }}
p-recipient3 = {{ .Name3 }}
p-recipient4 = {{ .Address.Street }}
p-recipient5 = {{ formatCity .Address.Zip .Address.City }}
p-recipient6 = {{ .Address.Country | toCIILabel }}
{{ end }}

# Info Block
{{ with .Invoice.Invoice.Buyer.Contact }}
p-infoContact = {{ .Name }}
p-infoPhone = {{ .Phone }}
p-infoFax = {{ .Fax }}
p-infoEmail = {{ .Email }}
{{ end }}

# gesetzliche Angaben
p-customerID = {{ .Invoice.Invoice.Seller.ID }}
p-invoiceID = {{ .Invoice.Invoice.InvoiceID }}
p-issueDate = {{ .Invoice.Invoice.IssueDate | toPARdate }}
p-serviceDate = {{ .Invoice.Invoice.ServiceDate | toPARdate }}

# Zahlungsbedingugen
{{ with .Invoice.Invoice.PaymentTerms }}
p-dueDescription = {{ .Description }}
p-dueDate = {{ .DueDate | toPARdate }}
{{ end }}

# Impressum
{{ with .Invoice.Footers }}
p-impressumAddresses = {{ .Addresses }}
p-impressumContacts = {{ .Contacts }}
p-impressumBank = {{ .Bank }}
p-impressumLegal = {{ .Legal }}
{{ end }}

# Tabellen Decimal Point = . Feld-Trenner = ;
# Präfix: t-
# lh: lineHeader; li: lineItem; ls: lineSum; vi: vatItem; vs: vatSum; zi: Payment amount Item
# Table is sorted by column 1 (t-leistung[aaxx])
# Tabelle mit Namen: leistung
{{ range .LineItems -}}
{{ . }}
{{ end }}


