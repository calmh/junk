#!/usr/bin/env R --vanilla --no-save -f

library(grid)
library(reshape2)
library(ggplot2)

d <- read.table("result.csv", sep=",", header=TRUE, colClasses=c("factor", "factor", "numeric"))
d$Month = as.Date(paste(as.character(d$Month), "-01", sep=""))
d$Value = d$Value / 100 / 1000
cats <- c("Nettoomsättning","Personalkostnader","Externa kostnader")
d$Category = factor(d$Category, levels=c("Income","Employees","Expenses"), labels=cats)

m <- data.frame(Category=cats)
m$Mean[m$Category == "Nettoomsättning"] = mean(d$Value[d$Category == "Nettoomsättning"])
m$Mean[m$Category == "Personalkostnader"] = mean(d$Value[d$Category == "Personalkostnader"])
m$Mean[m$Category == "Externa kostnader"] = mean(d$Value[d$Category == "Externa kostnader"])

p1 <- ggplot(d, aes(x=Month, y=Value)) +
    geom_hline(aes(yintercept=Mean), m) +
    geom_bar(stat="identity", width=15) +
    facet_wrap(~Category) +
    xlab(NULL) + ylab("t.kr")

d[d$Category == "Personalkostnader",]$Value = -d[d$Category == "Personalkostnader",]$Value
d[d$Category == "Externa kostnader",]$Value = -d[d$Category == "Externa kostnader",]$Value

r <- aggregate(x=d$Value, by=list(d$Month), FUN=sum)
r$Month = r$Group.1
r$Revenue = r$x
r$Positive = r$Revenue > 0
r$Accumulated = cumsum(r$x)
r$Facet = "Resultat"

p2 <- ggplot(r, aes(x=Month, y=Revenue)) +
    geom_bar(stat="identity", width=15, aes(fill = Positive)) +
    geom_point(aes(y=Accumulated)) +
    facet_wrap(~Facet) +
    xlab(NULL) + ylab("t.kr")

theme_set(theme_light(base_size = 12) + theme(legend.position="none"))
pdf("result.pdf", paper="a4r", width=0, height=0)
grid.newpage()
pushViewport(viewport(layout = grid.layout(2, 1, heights=unit(c(1,1), rep('null',2)))))
vplayout <- function(x, y)
    viewport(layout.pos.row=x, layout.pos.col=y)
print(p1, vp=vplayout(1, 1))
print(p2, vp=vplayout(2, 1))
dev.off()
