#!/usr/bin/env R --vanilla --no-save -f

library(grid)
library(reshape2)
library(ggplot2)

d <- read.table("issues.csv", sep=",", header=TRUE, colClasses=c("factor", "factor", "numeric"))

d <- d[d$commit == "true",]
d <- d[d$days < 200,]

p1 <- ggplot(d, aes(days)) +
    geom_bar() +
    xlab("Days") + ylab("Issues")


theme_set(theme_light(base_size = 12) + theme(legend.position="none"))
pdf("result.pdf", paper="a4r", width=0, height=0)
grid.newpage()
pushViewport(viewport(layout = grid.layout(2, 1, heights=unit(c(1,1), rep('null',2)))))
vplayout <- function(x, y)
    viewport(layout.pos.row=x, layout.pos.col=y)
print(p1, vp=vplayout(1, 1))
#print(p2, vp=vplayout(2, 1))
dev.off()
