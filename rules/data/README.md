# Event data

This data was obtained from the Provider dumps in
https://github.com/nasbench/EVTX-ETW-Resources/


Use the following VQL artifact

```yaml
name: ExtractEventSchema
type: SERVER

parameters:
   - name: PathToXML
     default: F:/evtx/EVTX-ETW-Resources/ETWProvidersManifests/Windows10/22H2/W10_22H2_Pro_20230321_19045.2728/WEPExplorer

sources:
  - query: |
        LET GetFieldNames(Template) = parse_xml(accessor="data", file=Template).template.data.Attrname

        SELECT *
        FROM foreach(
          row={
            SELECT
                   parse_xml(
                     file=OSPath) AS Parsed
            FROM glob(globs="Microsoft-Windows-*",
                      root=PathToXML)
            WHERE NOT Name =~ "All.xml"
          },
          query={
           SELECT Id, Channel, Message,
              if(condition=Template, then=GetFieldNames(Template=Template)) AS Fields
           FROM foreach(row=Parsed.Providers.Provider.EventMetadata.Event, column="_value")
           WHERE Channel
          })
```
