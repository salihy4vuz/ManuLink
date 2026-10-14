use super::shapes::{Part, Rect, Sheet};

pub enum FreeRectChoiceHeuristic {
    BottomLeftRule,
    BestShortSideFit,
}

pub struct MaxRectsPacker;

impl MaxRectsPacker {
    pub fn new() -> Self {
        Self {}
    }

    pub fn insert_part(&self, sheet: &mut Sheet, mut part: Part, method: FreeRectChoiceHeuristic) -> bool {
        let (best_node, best_score1, best_score2) = self.find_position_for_new_node(sheet, &part, method);

        if best_node.height == 0.0 {
            return false;
        }

        part.placed_x = best_node.x;
        part.placed_y = best_node.y;
        
        // Check if we swapped width and height
        if best_node.width != part.width {
            part.rotated = true;
        }

        let mut i = 0;
        while i < sheet.free_rects.len() {
            let free_node = sheet.free_rects[i].clone();
            if self.split_free_node(&mut sheet.free_rects, free_node, &best_node) {
                sheet.free_rects.remove(i);
                // i is not incremented because we removed an element
            } else {
                i += 1;
            }
        }
        self.prune_free_list(sheet);

        sheet.placed_parts.push(part);
        true
    }

    fn find_position_for_new_node(&self, sheet: &Sheet, part: &Part, method: FreeRectChoiceHeuristic) -> (Rect, f64, f64) {
        let mut best_node = Rect::new(0.0, 0.0, 0.0, 0.0);
        let mut best_short_side_fit = f64::MAX;
        let mut best_long_side_fit = f64::MAX;
        
        // Simple Bottom-Left Rule implementation for MVP
        let mut best_y = f64::MAX;
        let mut best_x = f64::MAX;

        for free_rect in &sheet.free_rects {
            // Try to place normally
            if free_rect.width >= part.width && free_rect.height >= part.height {
                match method {
                    FreeRectChoiceHeuristic::BottomLeftRule => {
                        let top_side_y = free_rect.y + part.height;
                        if top_side_y < best_y || (top_side_y == best_y && free_rect.x < best_x) {
                            best_node.x = free_rect.x;
                            best_node.y = free_rect.y;
                            best_node.width = part.width;
                            best_node.height = part.height;
                            best_y = top_side_y;
                            best_x = free_rect.x;
                        }
                    }
                    FreeRectChoiceHeuristic::BestShortSideFit => {
                        let leftover_horiz = (free_rect.width - part.width).abs();
                        let leftover_vert = (free_rect.height - part.height).abs();
                        let short_side_fit = leftover_horiz.min(leftover_vert);
                        let long_side_fit = leftover_horiz.max(leftover_vert);

                        if short_side_fit < best_short_side_fit || (short_side_fit == best_short_side_fit && long_side_fit < best_long_side_fit) {
                            best_node.x = free_rect.x;
                            best_node.y = free_rect.y;
                            best_node.width = part.width;
                            best_node.height = part.height;
                            best_short_side_fit = short_side_fit;
                            best_long_side_fit = long_side_fit;
                        }
                    }
                }
            }

            // Try to place rotated 90 degrees
            if free_rect.width >= part.height && free_rect.height >= part.width {
                match method {
                    FreeRectChoiceHeuristic::BottomLeftRule => {
                        let top_side_y = free_rect.y + part.width;
                        if top_side_y < best_y || (top_side_y == best_y && free_rect.x < best_x) {
                            best_node.x = free_rect.x;
                            best_node.y = free_rect.y;
                            best_node.width = part.height;
                            best_node.height = part.width;
                            best_y = top_side_y;
                            best_x = free_rect.x;
                        }
                    }
                    FreeRectChoiceHeuristic::BestShortSideFit => {
                        let leftover_horiz = (free_rect.width - part.height).abs();
                        let leftover_vert = (free_rect.height - part.width).abs();
                        let short_side_fit = leftover_horiz.min(leftover_vert);
                        let long_side_fit = leftover_horiz.max(leftover_vert);

                        if short_side_fit < best_short_side_fit || (short_side_fit == best_short_side_fit && long_side_fit < best_long_side_fit) {
                            best_node.x = free_rect.x;
                            best_node.y = free_rect.y;
                            best_node.width = part.height;
                            best_node.height = part.width;
                            best_short_side_fit = short_side_fit;
                            best_long_side_fit = long_side_fit;
                        }
                    }
                }
            }
        }
        (best_node, best_short_side_fit, best_long_side_fit)
    }

    fn split_free_node(&self, free_rects: &mut Vec<Rect>, mut free_node: Rect, used_node: &Rect) -> bool {
        // Test with SAT if the rectangles even intersect.
        if used_node.x >= free_node.x + free_node.width || used_node.x + used_node.width <= free_node.x ||
           used_node.y >= free_node.y + free_node.height || used_node.y + used_node.height <= free_node.y {
            return false;
        }

        if used_node.x < free_node.x + free_node.width && used_node.x + used_node.width > free_node.x {
            // New node at the top side of the used node.
            if used_node.y > free_node.y && used_node.y < free_node.y + free_node.height {
                let mut new_node = free_node.clone();
                new_node.height = used_node.y - new_node.y;
                free_rects.push(new_node);
            }

            // New node at the bottom side of the used node.
            if used_node.y + used_node.height < free_node.y + free_node.height {
                let mut new_node = free_node.clone();
                new_node.y = used_node.y + used_node.height;
                new_node.height = free_node.y + free_node.height - (used_node.y + used_node.height);
                free_rects.push(new_node);
            }
        }

        if used_node.y < free_node.y + free_node.height && used_node.y + used_node.height > free_node.y {
            // New node at the left side of the used node.
            if used_node.x > free_node.x && used_node.x < free_node.x + free_node.width {
                let mut new_node = free_node.clone();
                new_node.width = used_node.x - new_node.x;
                free_rects.push(new_node);
            }

            // New node at the right side of the used node.
            if used_node.x + used_node.width < free_node.x + free_node.width {
                let mut new_node = free_node.clone();
                new_node.x = used_node.x + used_node.width;
                new_node.width = free_node.x + free_node.width - (used_node.x + used_node.width);
                free_rects.push(new_node);
            }
        }

        true
    }

    fn prune_free_list(&self, sheet: &mut Sheet) {
        let mut i = 0;
        while i < sheet.free_rects.len() {
            let mut j = i + 1;
            while j < sheet.free_rects.len() {
                if self.is_contained_in(&sheet.free_rects[i], &sheet.free_rects[j]) {
                    sheet.free_rects.remove(i);
                    i -= 1; // Need to re-evaluate this index
                    break;
                }
                if self.is_contained_in(&sheet.free_rects[j], &sheet.free_rects[i]) {
                    sheet.free_rects.remove(j);
                    j -= 1;
                }
                j += 1;
            }
            i += 1;
        }
    }

    fn is_contained_in(&self, a: &Rect, b: &Rect) -> bool {
        a.x >= b.x && a.y >= b.y 
        && a.x + a.width <= b.x + b.width 
        && a.y + a.height <= b.y + b.height
    }
}
