<?xml version="1.0" encoding="UTF-8"?>
<rsm:CrossIndustryInvoice 
  xmlns:rsm="urn:un:unece:uncefact:data:standard:CrossIndustryInvoice:100" 
  xmlns:ram="urn:un:unece:uncefact:data:standard:ReusableAggregateBusinessInformationEntity:100" 
  xmlns:udt="urn:un:unece:uncefact:data:standard:UnqualifiedDataType:100">
<rsm:ExchangedDocumentContext>
  <ram:BusinessProcessSpecifiedDocumentContextParameter>
    <ram:ID>urn:fdc:peppol.eu:2017:poacc:billing:01:1.0</ram:ID>
  </ram:BusinessProcessSpecifiedDocumentContextParameter>
  <ram:GuidelineSpecifiedDocumentContextParameter>
    <ram:ID>urn:cen.eu:en16931:2017#compliant#urn:xeinkauf.de:kosit:xrechnung_3.0</ram:ID>
  </ram:GuidelineSpecifiedDocumentContextParameter>
</rsm:ExchangedDocumentContext>
<rsm:ExchangedDocument>
  <ram:ID>{{.InvoiceID}}</ram:ID>
  <ram:TypeCode>380</ram:TypeCode>
  <ram:IssueDateTime>
    <udt:DateTimeString format="102">{{.IssueDate  | toXMLdate }}</udt:DateTimeString>
  </ram:IssueDateTime>
</rsm:ExchangedDocument>
<rsm:SupplyChainTradeTransaction>
  {{- range .Lines}}
  <ram:IncludedSupplyChainTradeLineItem>
    <ram:AssociatedDocumentLineDocument>
      <ram:LineID>{{.Pos}}</ram:LineID>
    </ram:AssociatedDocumentLineDocument>
    <ram:SpecifiedTradeProduct>
      <ram:Name>{{.Description | xmlEscape }}</ram:Name>
    </ram:SpecifiedTradeProduct>
    <ram:SpecifiedLineTradeAgreement>
      <ram:NetPriceProductTradePrice>
        <ram:ChargeAmount>{{.UnitPrice | formatEuro }}</ram:ChargeAmount>
      </ram:NetPriceProductTradePrice>
    </ram:SpecifiedLineTradeAgreement>
    <ram:SpecifiedLineTradeDelivery>
      <ram:BilledQuantity unitCode="{{ .Unit | toCIIcode }}">{{ printf "%.2f" .Quantity}}</ram:BilledQuantity>
    </ram:SpecifiedLineTradeDelivery>
    <ram:SpecifiedLineTradeSettlement>
    <ram:ApplicableTradeTax>
        <ram:TypeCode>VAT</ram:TypeCode>
        <ram:CategoryCode>S</ram:CategoryCode> <ram:RateApplicablePercent>{{ .TaxRate | formatEuro }}</ram:RateApplicablePercent>
    </ram:ApplicableTradeTax>
      <ram:SpecifiedTradeSettlementLineMonetarySummation>
        <ram:LineTotalAmount>{{ .LineTotal | formatEuro }}</ram:LineTotalAmount>
      </ram:SpecifiedTradeSettlementLineMonetarySummation>
    </ram:SpecifiedLineTradeSettlement>
  </ram:IncludedSupplyChainTradeLineItem>
  {{- end }}
  <ram:ApplicableHeaderTradeAgreement>
    <ram:BuyerReference>{{ .BuyerReference | xmlEscape }}</ram:BuyerReference>
    <ram:SellerTradeParty>
      <ram:Name>{{.Seller.Name | xmlEscape }}</ram:Name>
      <ram:DefinedTradeContact>
        <ram:PersonName>{{ .Seller.Contact.Name }}</ram:PersonName>
        <ram:TelephoneUniversalCommunication>
            <ram:CompleteNumber>{{ .Seller.Contact.Phone }}</ram:CompleteNumber>
        </ram:TelephoneUniversalCommunication>
        <ram:EmailURIUniversalCommunication>
          <ram:URIID>{{ .Seller.Contact.Email }}</ram:URIID>
        </ram:EmailURIUniversalCommunication>
      </ram:DefinedTradeContact>
      <ram:PostalTradeAddress>
        <ram:PostcodeCode>{{.Seller.Address.Zip}}</ram:PostcodeCode>
        <ram:LineOne>{{.Seller.Address.Street | xmlEscape }}</ram:LineOne>
        <ram:CityName>{{.Seller.Address.City | xmlEscape }}</ram:CityName>
        <ram:CountryID>{{.Seller.Address.Country | toCIIcode}}</ram:CountryID>
      </ram:PostalTradeAddress>
      <ram:URIUniversalCommunication>
        <ram:URIID schemeID="EM">{{.Seller.Contact.Email}} </ram:URIID>
      </ram:URIUniversalCommunication>
      {{- if .Seller.TaxID}}
      <ram:SpecifiedTaxRegistration>
        <ram:ID schemeID="VA">{{.Seller.TaxID}}</ram:ID>
      </ram:SpecifiedTaxRegistration>
    {{end -}}
    </ram:SellerTradeParty>
    <ram:BuyerTradeParty>
      <ram:Name>{{.Buyer.Name}}</ram:Name>
      <ram:PostalTradeAddress>
      <ram:PostcodeCode>{{.Buyer.Address.Zip}}</ram:PostcodeCode>
        <ram:LineOne>{{.Buyer.Address.Street | xmlEscape }}</ram:LineOne>
        <ram:CityName>{{.Buyer.Address.City | xmlEscape }}</ram:CityName>
        <ram:CountryID>{{.Buyer.Address.Country | toCIIcode}}</ram:CountryID>
      </ram:PostalTradeAddress>
      {{- if .Buyer.Contact}}
      <ram:URIUniversalCommunication>
        <ram:URIID schemeID="EM">{{.Buyer.Contact.Email}} </ram:URIID>
      </ram:URIUniversalCommunication>
      {{end -}}
    </ram:BuyerTradeParty>
    </ram:ApplicableHeaderTradeAgreement>
    <ram:ApplicableHeaderTradeDelivery>
      <ram:ActualDeliverySupplyChainEvent>
        <ram:OccurrenceDateTime>
          <udt:DateTimeString format="102">{{ .ServiceDate | toXMLdate }}</udt:DateTimeString>
        </ram:OccurrenceDateTime>
      </ram:ActualDeliverySupplyChainEvent>
    </ram:ApplicableHeaderTradeDelivery>
    <ram:ApplicableHeaderTradeSettlement>
      <ram:InvoiceCurrencyCode>{{.Currency | toCIIcode }}</ram:InvoiceCurrencyCode>
      <ram:SpecifiedTradeSettlementPaymentMeans>
      <ram:TypeCode>{{ .PaymentMean | toCIIcode }}</ram:TypeCode>
      <ram:PayeePartyCreditorFinancialAccount>
        <ram:IBANID>{{ .Seller.IBAN }}</ram:IBANID>
      </ram:PayeePartyCreditorFinancialAccount>
      </ram:SpecifiedTradeSettlementPaymentMeans>
      {{ range .Vats -}}
      <ram:ApplicableTradeTax>
        <ram:CalculatedAmount>{{ .VatAmount | formatEuro}}</ram:CalculatedAmount>
        <ram:TypeCode>VAT</ram:TypeCode>
        <ram:BasisAmount>{{.BasisAmount | formatEuro}}</ram:BasisAmount>
        <ram:CategoryCode>{{ .Category }}</ram:CategoryCode>
        <ram:RateApplicablePercent>{{.TaxRate | formatEuro}}</ram:RateApplicablePercent>
      </ram:ApplicableTradeTax>
      {{end -}}
      {{- if or .ServicePeriod.Start .ServicePeriod.End -}}
      <ram:BillingSpecifiedPeriod>
      {{- if .ServicePeriod.Start -}}
        <ram:StartDateTime>
          <udt:DateTimeString format="102">{{.ServicePeriod.Start | toXMLdate}}</udt:DateTimeString>
        </ram:StartDateTime>
      {{ end -}}
      {{- if .ServicePeriod.End -}}
        <ram:EndDateTime>
          <udt:DateTimeString format="102">{{.ServicePeriod.End | toXMLdate }}</udt:DateTimeString>
        </ram:EndDateTime>
      {{end -}}
      </ram:BillingSpecifiedPeriod>
      {{end -}}
      {{- if or .PaymentTerms.Description .PaymentTerms.DueDate -}}
      <ram:SpecifiedTradePaymentTerms>
        {{- if .PaymentTerms.Description -}}
        <ram:Description>{{.PaymentTerms.Description}}</ram:Description>
        {{end -}}
        <ram:DueDateDateTime>
            <udt:DateTimeString format="102">{{.PaymentTerms.DueDate | toXMLdate }}</udt:DateTimeString>
        </ram:DueDateDateTime>
      </ram:SpecifiedTradePaymentTerms>      
      {{end -}}
      <ram:SpecifiedTradeSettlementHeaderMonetarySummation>
        <ram:LineTotalAmount>{{ .Totals.LineTotalAmount | formatEuro}}</ram:LineTotalAmount>
        {{ if ne .Totals.ChargeTotalAmount 0.0}}
        <ram:ChargeTotalAmount>{{ .Totals.ChargeTotalAmount | formatEuro}}</ram:ChargeTotalAmount>
        {{end -}}
        {{if ne .Totals.AllowanceTotalAmount 0.0}}
        <ram:AllowanceTotalAmount>{{ .Totals.AllowanceTotalAmount | formatEuro}}</ram:AllowanceTotalAmount>
        {{end -}}
        <ram:TaxBasisTotalAmount>{{ .Totals.LineTotalAmount | formatEuro }}</ram:TaxBasisTotalAmount>
        <ram:TaxTotalAmount currencyID="{{.Currency | toCIIcode }}">{{ .Totals.TaxTotalAmount | formatEuro}}</ram:TaxTotalAmount>
        <ram:GrandTotalAmount>{{ .Totals.GrandTotalAmount | formatEuro}}</ram:GrandTotalAmount>
        <ram:DuePayableAmount>{{ .Totals.DuePayableAmount | formatEuro}}</ram:DuePayableAmount>
      </ram:SpecifiedTradeSettlementHeaderMonetarySummation>
    </ram:ApplicableHeaderTradeSettlement>
  </rsm:SupplyChainTradeTransaction>
</rsm:CrossIndustryInvoice>
